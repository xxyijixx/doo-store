package app

import (
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/core/repo"
	"doo-store/backend/utils/docker"
	"fmt"
	"os"
	"path"

	"gorm.io/gorm"
)

func Init() {

	// 初始化目录
	constant.DataDir = resolveDataDir(config.EnvConfig.App().DATA_DIR)
	constant.AppInstallDir = path.Join(constant.DataDir, "apps")
	constant.NginxDir = path.Join(constant.DataDir, "nginx")

	fmt.Println("数据目录: ", constant.DataDir)
	fmt.Println("应用安装目录: ", constant.AppInstallDir)
	fmt.Println("Nginx配置目录: ", constant.NginxDir)

	plugins, err := repo.AppInstalled.Select(repo.AppInstalled.ID, repo.AppInstalled.IpAddress).Find()
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	usedIPs := []string{}
	for _, plugin := range plugins {
		usedIPs = append(usedIPs, plugin.IpAddress)
	}
	// 初始化全局IP分配器
	if err := docker.InitIPAllocator(config.EnvConfig.PLUGIN_CIDR, usedIPs); err != nil {
		fmt.Printf("初始化IP分配器失败: %v", err)
		panic(err)
	}

	// 加载默认数据
	LoadData()
}

func resolveDataDir(dataDir string) string {
	if dataDir != "" {
		return dataDir
	}

	var err error
	dataDir, err = os.Getwd()
	if err != nil {
		fmt.Printf("获取当前工作目录失败: %v\n", err)
		return ""
	}
	dataDir = path.Join(dataDir, "docker", "dood")

	return dataDir
}
