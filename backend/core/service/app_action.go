package service

import (
	"doo-store/backend/core/model"
	"doo-store/backend/core/repo"
	"doo-store/backend/utils/compose"
	"doo-store/backend/utils/docker"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PluginActinManager struct {
}

var pluginActionManager = PluginActinManager{}

// restart 重新启动插件
func (m PluginActinManager) Restart(appInstalled *model.AppInstalled, envContent string) error {
	appKey, composeFile := pluginHelper.GetAppKeyAndComposeFile(appInstalled.Key)
	_, err := compose.Down(composeFile)
	if err != nil {
		log.WithError(err).Error("执行docker compose down命令失败")
		return fmt.Errorf("执行docker compose down命令失败: %w", err)
	}
	_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, model.PluginStatusInstalling)
	// 写入docker-compose.yaml和环境文件
	composeFile, err = pluginHelper.WriteComposeFile(appKey, appInstalled.DockerCompose)
	if err != nil {
		log.Error("DockerCompose文件写入失败", err.Error())
		return err
	}
	_, err = pluginHelper.WriteEnvFile(appKey, envContent)
	if err != nil {
		log.Error("环境变量文件写入失败", err.Error())
		return err
	}
	stdout, err := compose.Up(composeFile)
	if err != nil {
		log.Error("执行docker compose up命令错误", stdout)
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, model.PluginStatusUpErr)
		return err
	}
	_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, model.PluginStatusRunning)

	return nil
}

// up 插件启动
func (m PluginActinManager) Up(appInstalled *model.AppInstalled, envContent string) error {
	appKey := pluginHelper.GetAppKey(appInstalled.Key)
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		composeFile, err := pluginHelper.WriteComposeFile(appKey, appInstalled.DockerCompose)
		if err != nil {
			log.Error("Error WriteFile", err)
			return err
		}
		_, err = pluginHelper.WriteEnvFile(appKey, envContent)
		if err != nil {
			log.Error("Error WriteFile", err)
			return err
		}
		stdout, err := compose.Up(composeFile)
		if err != nil {
			stdout, err = docker.ParseError(stdout, err)
			log.Error("Error docker compose up:", stdout, err)
			return err
		}
		// 执行一次docker compose ps更新状态
		containers, err := compose.ParseDockerComposePsOutput(composeFile)
		if err != nil {
			return err
		}
		for _, container := range containers {
			log.WithFields(log.Fields{
				// "container": container,
				"containerName":  container.Name,
				"containerState": container.State,
			}).Debug("Docker容器状态")
			_, _ = repo.Use(tx).AppServiceStatus.Where(repo.AppServiceStatus.InstallID.Eq(appInstalled.ID)).
				Where(repo.AppServiceStatus.ContainerName.Eq(container.Name)).Updates(
				model.AppServiceStatus{
					Status:  container.State,
					Message: "",
				},
			)
		}
		fmt.Println(stdout)
		_, err = repo.Use(tx).AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
			model.AppInstalled{
				Status:  model.PluginStatusRunning,
				Message: "",
			},
		)
		if err != nil {
			return err
		}
		return nil
	})
	stderr := ""
	if err != nil {
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
			model.AppInstalled{
				Status:  model.PluginStatusUpErr,
				Message: err.Error(),
			},
		)
		stderr = err.Error()
	}
	insertLog(appInstalled.ID, "插件启动", stderr)
	return err
}

func (m PluginActinManager) Start(appInstalled *model.AppInstalled) error {
	var (
		stdout string
		err    error
	)
	_, composeFile := pluginHelper.GetAppKeyAndComposeFile(appInstalled.Key)
	// 插件未正常启动，执行up操作
	if appInstalled.Status == model.PluginStatusUpErr {
		stdout, err = compose.Up(composeFile)
	} else {
		stdout, err = compose.Operate(composeFile, string(model.PluginActionStart))
	}
	if err != nil {
		log.Info("Error docker compose operate")

		_, err = docker.ParseError(stdout, err)
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
			map[string]interface{}{
				repo.AppInstalled.Message.ColumnName().String(): err.Error(),
			},
		)
		return err
	}
	_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
		map[string]interface{}{
			repo.AppInstalled.Status.ColumnName().String():  model.PluginStatusRunning,
			repo.AppInstalled.Message.ColumnName().String(): "",
		},
	)
	insertLog(appInstalled.ID, fmt.Sprintf("插件操作[%s]", "start"), stdout)
	return nil
}

// pluginActionStop 插件停止
func (m PluginActinManager) Stop(appInstalled *model.AppInstalled) error {
	_, composeFile := pluginHelper.GetAppKeyAndComposeFile(appInstalled.Key)
	_, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, model.PluginStatusStopped)
	if err != nil {
		return err
	}
	stdout, err := compose.Stop(composeFile)
	if err != nil {
		return fmt.Errorf("error docker compose stop: %s", err.Error())
	}
	insertLog(appInstalled.ID, "插件停止", stdout)
	return nil
}
