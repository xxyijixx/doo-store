package task

import (
	"context"
	"doo-store/backend/constant"
	"doo-store/backend/core/model"
	"doo-store/backend/core/repo"
	"doo-store/backend/utils/docker"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// DockerMonitor 容器监控任务
type DockerMonitor struct {
	client *client.Client
	ctx    context.Context
}

// NewDockerMonitor 创建新的Docker监控器
func NewDockerMonitor(ctx context.Context) (*DockerMonitor, error) {
	cli, err := docker.NewClient()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", constant.ErrDockerClientCreate, err)
	}

	return &DockerMonitor{
		client: cli.GetClient(),
		ctx:    ctx,
	}, nil
}

// StartMonitoring 开始监控任务
func (dm *DockerMonitor) StartMonitoring(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := dm.monitorContainers(); err != nil {
					log.Printf("Error monitoring containers: %v", err)
				}
			case <-dm.ctx.Done():
				return
			}
		}
	}()
}

// 监控容器状态
func (dm *DockerMonitor) monitorContainers() error {
	log.Debug("正在处理容器状态")

	apps, err := repo.AppInstalled.Find()
	if err != nil {
		return fmt.Errorf("%s: %v", constant.ErrDockerFindApps, err)
	}

	appStatusMap := dm.getContainerStatuses(apps)
	dm.updateAppStatuses(appStatusMap)
	log.Debug("结束处理容器状态")
	return nil
}

// 获取所有容器的状态
func (dm *DockerMonitor) getContainerStatuses(apps []*model.AppInstalled) map[string]string {
	filterArgs := filters.NewArgs()
	appStatusMap := make(map[string]string)

	// 将状态标记为 “init”表示这些容器是本系统的容器
	for _, app := range apps {
		filterArgs.Add("name", app.Name)

		appStatusMap[app.Name] = docker.CustomContainerStatusInit
	}

	containers, err := dm.client.ContainerList(dm.ctx, container.ListOptions{
		All:     true,
		Filters: filterArgs,
	})
	if err != nil {
		log.Errorf("Failed to list containers: %v", err)
		return appStatusMap
	}

	for _, container := range containers {
		containerName := strings.TrimPrefix(container.Names[0], "/")
		if appStatusMap[containerName] == docker.CustomContainerStatusInit {
			appStatusMap[containerName] = container.State
		}
	}
	log.Info("容器状态", appStatusMap)
	return appStatusMap
}

func (dm *DockerMonitor) updateAppStatus(appName string, status string, message string) {
	appInstalled, err := repo.AppInstalled.Select(repo.AppInstalled.ID, repo.AppInstalled.Status).Where(repo.AppInstalled.Name.Eq(appName)).First()
	if err != nil {
		log.Errorf("Failed to find app record for %s: %v", appName, err)
		return
	}

	// 跳过处理 Installing 状态的应用
	if strings.EqualFold(appInstalled.Status, model.PluginStatusInstalling) {
		log.Debugf("Skipping status update for app %s as it is in Installing state", appName)
		return
	}

	// 只有状态发生变化时才更新
	if appInstalled.Status != status && appInstalled.Status != model.PluginStatusUpErr {
		log.Debugf("更新状态 %s [%s]", status, message)
		_, err = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
			map[string]interface{}{
				repo.AppInstalled.Status.ColumnName().String():  status,
				repo.AppInstalled.Message.ColumnName().String(): message,
			},
		)
		if err != nil {
			log.Errorf("Failed to update app status for %s: %v", appName, err)
			return
		}
	}
}

// 根据容器状态更新应用状态
func (dm *DockerMonitor) updateAppStatuses(appStatusMap map[string]string) {
	for appName, status := range appStatusMap {
		switch status {
		case docker.ContainerStatusRunning:
			dm.updateAppStatus(appName, model.PluginStatusRunning, "")
		case docker.ContainerStatusExited:
			dm.handleExitedContainer(appName)
		case docker.ContainerStatusRestarting:
			dm.updateAppStatus(appName, model.PluginStatusRestarting, "Container is restarting")
		case docker.ContainerStatusPaused:
			dm.updateAppStatus(appName, model.PluginStatusPaused, "Container is paused")
		case docker.ContainerStatusDead:
			dm.updateAppStatus(appName, model.PluginStatusDead, "Container is in dead state")
		case docker.CustomContainerStatusInit:
			dm.updateAppStatus(appName, model.PluginStatusError, "Container is not existing")
		default:
			dm.updateAppStatus(appName, model.PluginStatusUnknown, fmt.Sprintf("Unknown state: %s", status))
		}
	}
}

// 处理退出的容器
func (dm *DockerMonitor) handleExitedContainer(appName string) {
	container, err := dm.client.ContainerInspect(dm.ctx, appName)
	if err != nil {
		log.Warnf("Failed to inspect container %s: %v", appName, err)
		return
	}
	if container.State.ExitCode == 0 {
		dm.updateAppStatus(appName, model.PluginStatusStopped, "Container stopped normally")
	} else {
		message := fmt.Sprintf("Container exited with code %d: %s", container.State.ExitCode, container.State.Error)
		dm.updateAppStatus(appName, model.PluginStatusError, message)
	}
}
