package constant

import (
	"doo-store/backend/config"
	"fmt"
	"os"
	"path"
)

var (
	DataDir            = getDataDir(config.EnvConfig.DATA_DIR)
	AppInstallDir      = path.Join(DataDir, "apps")
	NginxDir           = path.Join(DataDir, "nginx")
)

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
