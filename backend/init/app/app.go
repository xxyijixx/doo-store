package app

import (
	"context"
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/utils/docker"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	log "github.com/sirupsen/logrus"
)

func Init() {

	constant.DataDir = getDataDir(config.EnvConfig.DATA_DIR)
	constant.AppInstallDir = path.Join(constant.DataDir, "apps")
	constant.NginxDir = path.Join(constant.DataDir, "nginx")
	constant.NginxConfigDir = path.Join(constant.NginxDir, "conf.d")
	constant.NginxAppsConfigDir = path.Join(constant.NginxDir, "apps")

	if config.EnvConfig.ENV == "prod" {
		err := os.MkdirAll(constant.DataDir, 0755)
		if err != nil {
			fmt.Println("创建目录失败")
			return
		}
	}

	dirs := []string{constant.DataDir, constant.AppInstallDir, constant.NginxDir}

	constant.AppNetwork = config.EnvConfig.APP_PREFIX + "network"

	for _, dir := range dirs {
		createDir(dir)
	}

	LoadData()

	_ = docker.CreateDefaultDockerNetwork()

	InitNginxProxy()

}

func getDataDir(dataDir string) string {
	if dataDir == "" {
		var err error
		dataDir, err = os.Getwd()
		if err != nil {
			fmt.Printf("获取当前工作目录失败: %v\n", err)
			return ""
		}
		dataDir = path.Join(dataDir, "docker", "dood")
	}
	return dataDir
}

func createDir(dirPath string) {
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		if os.IsExist(err) {
			return
		}
	}
}

func InitNginxProxy() {
	ImageNginxName := "nginx:alpine"
	client, err := docker.NewDockerClient()
	if err != nil {
		return
	}

	ctx := context.Background()

	nginxImage, err := client.ImageList(ctx, image.ListOptions{
		Filters: filters.NewArgs(
			filters.Arg("reference", ImageNginxName),
		),
	})
	if err != nil {
		log.Debug("获取镜像列表失败", err)
		return
	}
	if len(nginxImage) == 0 {
		reader, err := client.ImagePull(ctx, ImageNginxName, image.PullOptions{})
		if err != nil {
			log.Debug("拉取镜像失败", err)
			return
		}
		io.Copy(os.Stdout, reader)
	} else {
		log.Debugf("镜像%s已存在", ImageNginxName)
	}

	// 定义容器配置
	containerConfig := &container.Config{
		Image: ImageNginxName,
	}

	// 定义主机配置
	hostConfig := &container.HostConfig{
		ShmSize: 2 * 1024 * 1024 * 1024, // 2 GB
		Resources: container.Resources{
			Ulimits: []*units.Ulimit{
				{
					Name: "core",
					Hard: 0,
					Soft: 0,
				},
			},
		},
		// 挂载目录

		Binds: []string{
			fmt.Sprintf("%s:/etc/nginx/conf.d", constant.NginxDir),
		},
		// 重启策略
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8081",
				},
			},
		},
	}

	endpointsConfig := map[string]*network.EndpointSettings{
		constant.AppNetwork: {},
	}
	// 添加外部网络
	if config.EnvConfig.EXTERNAL_NETWORK_NAME != "" {
		endpointsConfig[config.EnvConfig.EXTERNAL_NETWORK_NAME] = &network.EndpointSettings{
			Aliases:   []string{constant.NginxContainerName}, // 添加别名
			IPAddress: config.EnvConfig.EXTERNAL_NETWORK_IP,
			Gateway:   config.EnvConfig.EXTERNAL_NETWORK_GATEWAY,
			IPAMConfig: &network.EndpointIPAMConfig{
				IPv4Address: config.EnvConfig.EXTERNAL_NETWORK_IP,
			},
		}
	}
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: endpointsConfig,
	}

	resp, err := client.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, constant.NginxContainerName)
	if err != nil {
		fmt.Println("创建容器失败", err)
		return
	}

	log.WithField("container_id", resp.ID).Debug("nginx容器创建成功")
	err = client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Debug("Nginx容器启动失败", err)
		return
	}
	log.WithField("container_id", resp.ID).Debug("nginx容器启动成功")
}
