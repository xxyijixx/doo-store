package compose

import (
	"doo-store/backend/config"
	"doo-store/backend/utils/cmd"
	"fmt"
)

func Pull(filePath string) (string, error) {
	insertPart := getInsertPart()
	stdout, err := cmd.Execf("docker-compose%s -f %s pull", insertPart, filePath)
	return stdout, err
}

func Up(filePath string) (string, error) {
	insertPart := getInsertPart()
	stdout, err := cmd.Execf("docker-compose%s -f %s up -d", insertPart, filePath)
	return stdout, err
}

func Down(filePath string) (string, error) {
	insertPart := getInsertPart()
	stdout, err := cmd.Execf("docker-compose%s -f %s down", insertPart, filePath)
	return stdout, err
}

func Start(filePath string) (string, error) {
	insertPart := getInsertPart()
	stdout, err := cmd.Execf("docker-compose%s -f %s start", insertPart, filePath)
	return stdout, err
}

func Stop(filePath string) (string, error) {
	insertPart := getInsertPart()
	stdout, err := cmd.Execf("docker-compose%s -f %s stop", insertPart, filePath)
	return stdout, err
}

func Restart(filePath string) (string, error) {
	insertPart := getInsertPart()
	stdout, err := cmd.Execf("docker-compose%s -f %s restart", insertPart, filePath)
	return stdout, err
}

func Operate(filePath, operation string) (string, error) {
	insertPart := getInsertPart()
	stdout, err := cmd.Execf("docker-compose%s -f %s %s", insertPart, filePath, operation)
	return stdout, err
}

func getInsertPart() string {
	if config.EnvConfig.App().SHARED_COMPOSE {
		return fmt.Sprintf(" -p %s", config.EnvConfig.App().SHARED_COMPOSE_NAME)
	}
	return ""
}
