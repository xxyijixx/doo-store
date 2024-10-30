package nginx

import (
	"context"
	"doo-store/backend/constant"
	"doo-store/backend/utils/docker"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	log "github.com/sirupsen/logrus"
)

func imageExposedPort(imageName string) (int, error) {
	ctx := context.Background()
	client, err := docker.NewDockerClient()
	if err != nil {
		return 0, err
	}
	imageInspect, _, err := client.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		log.Debugf("Failed to inspect image: %v", err)
		return 0, err
	}
	portNum := 0
	for port := range imageInspect.Config.ExposedPorts {
		portStr := strings.Split(string(port), "/")[0]
		portNum, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Failed to convert port to number: %v", err)
		}
		fmt.Printf("First exposed port (as number): %d\n", portNum)
		break
	}
	return portNum, nil
}

// AddLocation 添加一个location块
func AddLocation(locationName, proxyServerName string, port int) {
	locationPath := fmt.Sprintf("%s/%s.conf", constant.NginxAppsConfigDir, locationName)

	fileInfo, err := os.Stat(locationPath)
	if err != nil && !os.IsNotExist(err) {
		log.Debug("写入文件失败", err, fileInfo)
		return
	}

	fileContent := fmt.Sprintf(`location /%s/ {
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

	err = os.WriteFile(locationPath, []byte(fileContent), 0644)
	if err != nil {
		log.Debug("写入文件失败")
		panic(err)
	}

	err = reloadNginx()
	if err != nil {
		log.Debug("Nginx 重载失败", err)
	}
}

// 重载Nginx
func reloadNginx() error {
	client, err := docker.NewClient()
	if err != nil {
		log.Debug("获取Docker客户端失败")
		return err
	}
	list, err := client.ListContainersByName([]string{"nginx-core-proxy"})
	if err != nil {
		log.Debug("查找容器失败", err)
		return err
	}
	if len(list) < 1 {
		log.WithField("container_name", "ContainerNginxName").Debug("Nginx 容器不存在")
		return fmt.Errorf("nginx container not found")
	}

	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"nginx", "-s", "reload"},
	}

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return err
	}
	execIDResp, err := dockerClient.ContainerExecCreate(context.Background(), list[0].ID, execConfig)
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
	} else {
		fmt.Println("Nginx configuration reloaded successfully.")
	}

	return nil
}
