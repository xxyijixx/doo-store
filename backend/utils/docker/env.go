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
		envJson = string(jsonData)
	}
	return
}
