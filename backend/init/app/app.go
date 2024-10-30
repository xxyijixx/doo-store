package app

import (
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"fmt"
	"os"
	"path"
)

func Init() {
	constant.DataDir = getDataDir(config.EnvConfig.DATA_DIR)
	constant.AppInstallDir = path.Join(constant.DataDir, "apps")

	if config.EnvConfig.ENV == "prod" {
		err := os.MkdirAll(constant.DataDir, 0755)
		if err != nil {
			fmt.Println("创建目录失败")
			return
		}
	}

	dirs := []string{constant.DataDir, constant.AppInstallDir}

	for _, dir := range dirs {
		createDir(dir)
	}

	// _ = docker.CreateDefaultDockerNetwork()

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
