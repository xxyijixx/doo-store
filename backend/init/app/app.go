package app

import (
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/core/repo"
	"doo-store/backend/utils/docker" // add this line
	"fmt"                            // add this line
	"os"
	"path"

	"gorm.io/gorm"
)

func Init() {

	constant.DataDir = resolveDataDir(config.EnvConfig.DATA_DIR)
	constant.AppInstallDir = path.Join(constant.DataDir, "apps")
	constant.NginxDir = path.Join(constant.DataDir, "nginx")
	constant.NginxConfigDir = path.Join(constant.NginxDir, "conf.d")
	constant.NginxAppsConfigDir = path.Join(constant.NginxDir, "apps")

	if config.EnvConfig.ENV == "prod" {
		constant.DooTaskUrl = "http://" + config.EnvConfig.DOOTASK_APP_IPPR + ".3"
		if err := ensureDir(constant.DataDir); err != nil {
			fmt.Printf("创建目录失败: %v\n", err)
			return
		}
	}

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

	requiredDirs := []string{constant.DataDir, constant.AppInstallDir, constant.NginxDir}

	for _, dir := range requiredDirs {
		if err := ensureDir(dir); err != nil {
			fmt.Printf("创建目录失败: %v\n", err)
			return
		}
	}

	// 加载默认数据
	LoadData()
}

func resolveDataDir(dataDir string) string {
	if dataDir != "" {
		return dataDir
	}

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前工作目录失败: %v\n", err)
		return ""
	}
	return path.Join(workingDir, "docker", "dood")
}

func ensureDir(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
