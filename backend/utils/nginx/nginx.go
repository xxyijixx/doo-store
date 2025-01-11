package nginx

import (
	"bytes"
	"context"
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/utils/docker"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"text/template"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	log "github.com/sirupsen/logrus"
)

// AddLocation 添加一个location块
func AddLocation(tmpl, locationName, proxyServerName string, port int) error {
	locationPath := fmt.Sprintf("%s/%s.conf", constant.NginxAppsConfigDir, locationName)

	fileInfo, err := os.Stat(locationPath)
	if err != nil && !os.IsNotExist(err) {
		log.Debug("写入文件失败", err, fileInfo)
		return errors.New(constant.ErrNginxWriteFile)
	}
	fileContent := tmpl
	// 如果模板为空，使用默认配置
	if tmpl == "" {
		proxyPass := fmt.Sprintf("http://%s/", proxyServerName)
		if port != 0 {
			proxyPass = fmt.Sprintf("http://%s:%d/", proxyServerName, port)
		}
		fileContent = fmt.Sprintf(`location /plugin/%s/ {
	proxy_http_version 1.1;
	proxy_set_header X-Real-IP $remote_addr;
	proxy_set_header X-Real-PORT $remote_port;
	proxy_set_header X-Forwarded-Host $the_host;
	proxy_set_header X-Forwarded-Proto $the_scheme;
	proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	proxy_set_header Host $http_host;
	proxy_set_header Scheme $scheme;
	proxy_set_header Server-Protocol $server_protocol;
	proxy_set_header Server-Name $server_name;
	proxy_set_header Server-Addr $server_addr;
	proxy_set_header Server-Port $server_port;
	proxy_set_header Upgrade $http_upgrade;
	proxy_set_header Connection $connection_upgrade;
	proxy_read_timeout 3600s;
	proxy_send_timeout 3600s;
	proxy_connect_timeout 3600s;
	proxy_pass %s;
}`, locationName, proxyPass)
	} else {
		t, err := template.New("nginx").Parse(tmpl)
		if err != nil {
			log.Debug("解析模板内容失败:", err)
			return errors.New(constant.ErrNginxParseContent)
		}
		var buf bytes.Buffer
		t.Execute(&buf, map[string]interface{}{
			"Key":           locationName,
			"ContainerName": proxyServerName,
			"Port":          port,
		})

		fileContent = buf.String()
	}

	err = os.WriteFile(locationPath, []byte(fileContent), 0644)
	if err != nil {
		log.Debug("写入文件失败")
		return errors.New(constant.ErrNginxWriteFile)
	}

	nginxContainer, err := getNginxContainer()
	if err != nil {
		log.Debug("获取Nginx容器失败", err)
		return errors.New(constant.ErrNginxGetContainer)
	}
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Debug("创建Docker客户端失败", err)
		return errors.New(err.Error())
	}

	// 检查是否存在默认配置文件
	defaultConfPath := fmt.Sprintf("/etc/nginx/conf.d/apps/%s-default.conf", locationName)
	exists, err := dockerClient.FileExistsInContainer(nginxContainer.ID, defaultConfPath)
	if err != nil {
		log.Debug("检查默认配置文件失败", err)
		return errors.New(err.Error())
	}

	if exists {
		// 如果存在默认配置,重命名为.bak
		err = dockerClient.MoveFileWithCheck(nginxContainer.ID, defaultConfPath, defaultConfPath+".bak")
		if err != nil {
			log.Debug("重命名默认配置文件失败", err)
			return errors.New(err.Error())
		}
	}

	// 复制新的配置文件到容器
	err = dockerClient.CopyFileToContainer(nginxContainer.ID, locationPath, fmt.Sprintf("/etc/nginx/conf.d/apps/%s.conf", locationName))
	if err != nil {
		log.Debug("复制文件到容器失败", err)
		// 如果之前重命名了默认配置,需要恢复
		if exists {
			_ = dockerClient.MoveFileInContainer(nginxContainer.ID, defaultConfPath+".bak", defaultConfPath)
		}
		return errors.New(err.Error())
	}

	err = testNginxConfig(dockerClient.GetClient(), nginxContainer.ID)
	if err != nil {
		// 检测失败需要移除配置文件并恢复默认配置
		log.Info("Nginx 配置未通过检测", err)
		_ = dockerClient.RemoveFileFormContainer(nginxContainer.ID, fmt.Sprintf("/etc/nginx/conf.d/apps/%s.conf", locationName))
		if exists {
			_ = dockerClient.MoveFileInContainer(nginxContainer.ID, defaultConfPath+".bak", defaultConfPath)
		}
		return err
	}

	err = reloadNginx()
	if err != nil {
		log.Debug("Nginx 重载失败", err)
		return err
	}
	return nil
}

func RemoveLocation(locationName string) error {
	locationPath := fmt.Sprintf("%s/%s.conf", constant.NginxAppsConfigDir, locationName)

	nginxContainer, err := getNginxContainer()
	if err != nil {
		return err
	}
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}

	// 检查是否存在.bak文件
	bakPath := fmt.Sprintf("/etc/nginx/conf.d/apps/%s-default.conf.bak", locationName)
	exists, err := dockerClient.FileExistsInContainer(nginxContainer.ID, bakPath)
	if err != nil {
		log.Debug("检查.bak文件失败", err)
		return err
	}

	// 删除当前配置文件
	err = dockerClient.RemoveFileFormContainer(nginxContainer.ID, fmt.Sprintf("/etc/nginx/conf.d/apps/%s.conf", locationName))
	if err != nil {
		log.Debug("从容器中删除文件失败", err)
		return err
	}

	// 如果存在.bak文件，恢复它
	if exists {
		err = dockerClient.MoveFileInContainer(nginxContainer.ID, bakPath, fmt.Sprintf("/etc/nginx/conf.d/apps/%s-default.conf", locationName))
		if err != nil {
			log.Debug("恢复.bak文件失败", err)
			return err
		}
	}

	err = os.Remove(locationPath)
	if err != nil {
		fmt.Printf("删除文件失败: %v\n", err)
		return err
	}

	err = reloadNginx()
	if err != nil {
		return err
	}
	return nil
}

func ExtractLocations(nginxConfig string) []string {
	// 定义正则表达式来匹配 location 块
	re := regexp.MustCompile(`location\s+(/[^/]+(?:/[^/]+)*/*)\s+{`)
	matches := re.FindAllStringSubmatch(nginxConfig, -1)

	// 提取匹配的地址
	var locations []string
	for _, match := range matches {
		if len(match) > 1 {
			locations = append(locations, match[1])
		}
	}
	return locations
}

func getNginxContainer() (types.Container, error) {
	client, err := docker.NewClient()
	if err != nil {
		log.Debug("获取Docker客户端失败", err.Error())
		return types.Container{}, errors.New(constant.ErrNginxGetContainer)
	}

	list, err := client.ListContainersByName([]string{config.EnvConfig.GetNginxContainerName()})
	if err != nil {
		log.Debug("查找容器失败", err)
		return types.Container{}, errors.New(constant.ErrNginxGetContainer)
	}
	if len(list) < 1 {
		log.WithField("container_name", config.EnvConfig.GetNginxContainerName()).Debug("Nginx 容器不存在")
		return types.Container{}, errors.New(constant.ErrNginxContainerNotFound)
	}

	nginxContainer := list[0]
	return nginxContainer, nil
}

// 重载Nginx
func reloadNginx() error {

	nginxContainer, err := getNginxContainer()
	if err != nil {
		log.Info("获取Nginx容器失败")
		return err
	}
	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		log.Info("获取Docker Client失败", err)
		return err
	}
	err = testNginxConfig(dockerClient, nginxContainer.ID)
	if err != nil {
		log.Info("Nginx 配置未通过检测", err)
		return err
	}

	err = reloadNginxConfig(dockerClient, nginxContainer.ID)

	return err
}

func reloadNginxConfig(dockerClient *client.Client, containerID string) error {
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"nginx", "-s", "reload"},
	}

	execIDResp, err := dockerClient.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		return fmt.Errorf("error creating exec: %w", err)
	}

	execAttachResp, err := dockerClient.ContainerExecAttach(context.Background(), execIDResp.ID, container.ExecStartOptions{})
	if err != nil {
		panic(err)
	}
	defer execAttachResp.Close()

	outputDone := make(chan error)
	go func() {
		_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, execAttachResp.Reader)
		outputDone <- err
	}()

	err = <-outputDone
	if err != nil && err != io.EOF {
		fmt.Printf("Error during command execution: %v\n", err)
		return err
	} else {
		fmt.Println("Nginx configuration reloaded successfully.")
	}
	return nil
}

func testNginxConfig(dockerClient *client.Client, containerID string) error {

	// 创建一个执行命令的配置
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"nginx", "-t"},
	}

	// 创建执行命令
	execIDResp, err := dockerClient.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		return fmt.Errorf("error creating exec: %v", err)
	}

	// 执行命令
	execAttachResp, err := dockerClient.ContainerExecAttach(context.Background(), execIDResp.ID, container.ExecStartOptions{})
	if err != nil {
		return fmt.Errorf("error attaching to exec: %v", err)
	}
	defer execAttachResp.Close()

	// 读取命令输出
	outputDone := make(chan error)
	go func() {
		_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, execAttachResp.Reader)
		outputDone <- err
	}()

	// 等待命令执行完成
	err = <-outputDone
	if err != nil && err != io.EOF {
		return fmt.Errorf("error during command execution: %v", err)
	}

	return nil
}
