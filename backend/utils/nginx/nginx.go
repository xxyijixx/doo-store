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
		return errors.New("写入文件失败")
	}
	fileContent := tmpl
	// 如果模板为空，使用默认配置
	if tmpl == "" {
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
	proxy_pass http://%s:%d/;
}`, locationName, proxyServerName, port)
	} else {
		t, err := template.New("nginx").Parse(tmpl)
		if err != nil {
			log.Debug("解析模板内容失败:", err)
			return errors.New("解析内容失败")
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
		return errors.New("写入文件失败")
	}

	nginxContainer, err := getNginxContainer()
	if err != nil {
		log.Debug("获取Nginx容器失败", err)
		return errors.New("获取Nginx容器失败")
	}
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Debug("创建Docker客户端失败", err)
		return errors.New(err.Error())
	}

	err = dockerClient.CopyFileToContainer(nginxContainer.ID, locationPath, fmt.Sprintf("/etc/nginx/conf.d/site/%s.conf", locationName))
	if err != nil {
		log.Debug("复制文件到容器失败", err)
		return errors.New(err.Error())
	}

	err = testNginxConfig(dockerClient.GetClient(), nginxContainer.ID)
	if err != nil {
		// 检测失败需要移除配置文件
		log.Info("Nginx 配置未通过检测", err)
		err = dockerClient.RemoveFileFormContainer(nginxContainer.ID, fmt.Sprintf("/etc/nginx/conf.d/site/%s.conf", locationName))
		if err != nil {
			log.Debug("从容器中删除文件失败", err)
			return err
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
	err = dockerClient.RemoveFileFormContainer(nginxContainer.ID, fmt.Sprintf("/etc/nginx/conf.d/site/%s.conf", locationName))
	if err != nil {
		log.Debug("从容器中删除文件失败", err)
		return err
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

func getNginxContainer() (types.Container, error) {
	client, err := docker.NewClient()
	if err != nil {
		log.Debug("获取Docker客户端失败", err.Error())
		return types.Container{}, err
	}

	list, err := client.ListContainersByName([]string{config.EnvConfig.GetNginxContainerName()})
	if err != nil {
		log.Debug("查找容器失败", err)
		return types.Container{}, err
	}
	if len(list) < 1 {
		log.WithField("container_name", config.EnvConfig.GetNginxContainerName()).Debug("Nginx 容器不存在")
		return types.Container{}, fmt.Errorf("nginx container not found")
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
