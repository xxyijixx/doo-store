package service

import (
	"doo-store/backend/config"
	"doo-store/backend/constant"
	schemasReq "doo-store/backend/core/schemas/req"
	"doo-store/backend/utils/compose"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type PluginHelper struct {
}

var pluginHelper = PluginHelper{}

func (h PluginHelper) GetAppKey(key string) string {
	return config.EnvConfig.App().APP_PREFIX + key
}

func (h PluginHelper) GetComposeFile(key string) string {
	appKey := h.GetAppKey(key)
	composeFile := fmt.Sprintf("%s/%s/docker-compose.yml", constant.AppInstallDir, appKey)
	
	return composeFile
}

func (h PluginHelper) GetAppKeyAndComposeFile(key string) (appKey, composeFile string) {
	appKey = h.GetAppKey(key)
	composeFile = h.GetComposeFile(appKey)
	return
}

func (h PluginHelper) GenEnv(genEnvReq schemasReq.GenEnvReq) (envContent, envJson string, err error) {
	envFile := fmt.Sprintf("%s/%s/.env", constant.AppInstallDir, genEnvReq.AppKey)
	envContent = fmt.Sprintf("%s=%s\n", "CONTAINER_NAME", genEnvReq.ContainerName)
	envContent += fmt.Sprintf("%s=%s\n", "IP_ADDRESS", genEnvReq.IPAddress)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DIR", config.EnvConfig.DooTask().DIR)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_APP_ID", config.EnvConfig.DooTask().APP_ID)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_APP_IPPR", config.EnvConfig.DooTask().APP_IPPR)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_APP_KEY", config.EnvConfig.DooTask().APP_KEY)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_NETWORK_NAME", config.EnvConfig.App().NETWORK_NAME)

	// 数据库相关配置
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_HOST", config.EnvConfig.DooTaskDB().HOST)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_PORT", config.EnvConfig.DooTaskDB().PORT)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_DATABASE", config.EnvConfig.DooTaskDB().DATABASE)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_USERNAME", config.EnvConfig.DooTaskDB().USERNAME)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_PASSWORD", config.EnvConfig.DooTaskDB().PASSWORD)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_DB_PREFIX", config.EnvConfig.DooTaskDB().PREFIX)

	// Redis相关
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_REDIS_HOST", config.EnvConfig.DooTaskRedis().HOST)
	envContent += fmt.Sprintf("%s=%s\n", "DOOTASK_REDIS_PORT", config.EnvConfig.DooTaskRedis().PORT)

	for key, value := range genEnvReq.Envs {
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
	if genEnvReq.WriteFile {
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
func (h PluginHelper) WriteEnvFile(appKey, envContent string) (string, error) {
	envFile := h.GetEnvFile(appKey)
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write environment file: %w", err)
	}
	return envFile, nil
}

// 写Compose File
func (h PluginHelper) WriteComposeFile(appKey, composeContent string) (string, error) {
	composeFile := h.GetComposeFile(appKey)
	// 替换部分环境变量
	composeContent = compose.ReplaceEnvVariables(composeContent)
	err := os.WriteFile(composeFile, []byte(composeContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write docker-compose file: %w", err)
	}
	return composeFile, nil
}

func (h PluginHelper) GetEnvFile(appKey string) string {
	envFile := fmt.Sprintf("%s/%s/.env", constant.AppInstallDir, appKey)
	return envFile
}
