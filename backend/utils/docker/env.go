package docker

import (
	"doo-store/backend/constant"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GenEnv(appKey, containerName string, envs map[string]any, writeFile bool) (envContent, envJson string, err error) {
	envFile := fmt.Sprintf("%s/%s/.env", constant.AppInstallDir, appKey)
	envContent = fmt.Sprintf("%s=%s\n", "CONTAINER_NAME", containerName)
	for key, value := range envs {
		envContent += fmt.Sprintf("%s=%s\n", key, value)
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
func WrietEnvFile(appKey, envContent string) (string, error) {
	envFile := GetEnvFile(appKey)
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		log.Debug("Error WriteFile", err)
		return "", err
	}
	return envFile, nil
}

// 写Compose文件
func WriteComposeFile(appKey, composeContent string) (string, error) {
	composeFile := GetComposeFile(appKey)
	err := os.WriteFile(composeFile, []byte(composeContent), 0644)
	if err != nil {
		log.Debug("Error WriteFile", err)
		return "", err
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
