package docker

import (
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/utils/compose"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GenEnv(appKey, containerName, ipAddress string, envs map[string]any, writeFile bool) (envContent, envJson string, err error) {
	envFile := fmt.Sprintf("%s/%s/.env", constant.AppInstallDir, appKey)
	envContent = fmt.Sprintf("%s=%s\n", "CONTAINER_NAME", containerName)
	envContent += fmt.Sprintf("%s=%s\n", "IP_ADDRESS", ipAddress)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DIR", config.EnvConfig.GetDootaskDir())
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_APP_ID", config.EnvConfig.DOOTASK_APP_ID)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_APP_IPPR", config.EnvConfig.DOOTASK_APP_IPPR)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_APP_KEY", config.EnvConfig.DOOTASK_APP_KEY)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_NETWORK_NAME", config.EnvConfig.GetNetworkName())

	// 数据库相关配置
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_HOST", config.EnvConfig.DOOTASK_DB_HOST)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_PORT", config.EnvConfig.DOOTASK_DB_PORT)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_DATABASE", config.EnvConfig.DOOTASK_DB_DATABASE)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_USERNAME", config.EnvConfig.DOOTASK_DB_USERNAME)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_PASSWORD", config.EnvConfig.DOOTASK_DB_PASSWORD)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_PREFIX", config.EnvConfig.DOOTASK_DB_PREFIX)

	for key, value := range envs {
		var envValue string
		switch v := value.(type) {
		case float64:
			envValue = fmt.Sprintf("%f", v)
		case string:
			envValue = v
		default:
			envValue = fmt.Sprintf("%v", v)
		}
		envContent += fmt.Sprintf("%s=%s\n", key, envValue)
	}
	if writeFile {
		err = os.WriteFile(envFile, []byte(envContent), 0644)
		if err != nil {
			log.Debug("Error WriteFile", err)
			return
		}
	}
	envMap := map[string]string{}
	envContentLine := strings.Split(envContent, "\n")
	for _, line := range envContentLine {
		env := strings.Split(line, "=")
		if len(env) != 2 {
			continue
		}
		envMap[env[0]] = env[1]
	}
	jsonData, err := json.Marshal(envMap)
	if err != nil {
		return
	}
	envJson = string(jsonData)
	return
}

// 写环境变量文件
func WriteEnvFile(appKey, envContent string) (string, error) {
	envFile := GetEnvFile(appKey)
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write environment file: %w", err)
	}
	return envFile, nil
}

// 写Compose文件
func WriteComposeFile(appKey, composeContent string) (string, error) {
	composeFile := GetComposeFile(appKey)
	// 替换部分环境变量
	composeContent = compose.ReplaceEnvVariables(composeContent)
	err := os.WriteFile(composeFile, []byte(composeContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write docker-compose file: %w", err)
	}
	return composeFile, nil
}

func GetComposeFile(appKey string) string {
	composeFile := fmt.Sprintf("%s/%s/docker-compose.yml", constant.AppInstallDir, appKey)
	return composeFile
}

func GetEnvFile(appKey string) string {
	envFile := fmt.Sprintf("%s/%s/.env", constant.AppInstallDir, appKey)
	return envFile
}
