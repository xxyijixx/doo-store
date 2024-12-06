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
func NewDockerMonitor() (*DockerMonitor, error) {
	cli, err := docker.NewClient()
	if err != nil {
		return nil, fmt.Errorf(constant.ErrDockerClientCreate, err)
	}

	return &DockerMonitor{
		client: cli.GetClient(),
		ctx:    context.Background(),
	}, nil
}

// StartMonitoring 开始监控任务
func (dm *DockerMonitor) StartMonitoring(interval time.Duration) {
	task := func() error {
		log.Info("容器状态监听")
		apps, err := repo.AppInstalled.Find()
		if err != nil {
			return fmt.Errorf(constant.ErrDockerFindApps, err)
		}

		// Create filter args with app names
		filterArgs := filters.NewArgs()
		appNameMap := make(map[string]bool)
		appStatusMap := make(map[string]string) // Track container status for each app
		for _, app := range apps {
			filterArgs.Add("name", app.Name)
			appNameMap[app.Name] = true
			appStatusMap[app.Name] = "unknown" // Initialize status
		}

		// List all containers without filtering
		containers, err := dm.client.ContainerList(dm.ctx, container.ListOptions{
			All:     true,
			Filters: filterArgs,
		})
		if err != nil {
			return fmt.Errorf(constant.ErrDockerListContainers, err)
		}

		// Process container statuses
		for _, container := range containers {
			// fmt.Printf("Container ID: %s\n", container.ID[:12])
			// fmt.Printf("Image: %s\n", container.Image)
			// fmt.Printf("Status: %s\n", container.Status)
			// fmt.Printf("State: %s\n", container.State)
			// fmt.Printf("Names: %v\n", container.Names)

			// Check if container name exists in apps
			containerName := container.Names[0][1:] // Remove leading slash

			state := container.State
			appStatusMap[containerName] = state // Update status in map

			switch state {
			case "exited":
				// Get container details to check exit code
				inspect, err := dm.client.ContainerInspect(dm.ctx, container.ID)
				if err != nil {
					log.Errorf("Failed to inspect container %s: %v", containerName, err)
					continue
				}
				exitCode := inspect.State.ExitCode
				log.Errorf("Container %s exited with code %d: %s", containerName, exitCode, inspect.State.Error)

			case "running":
				// Container is running normally
				log.Debugf("Container %s is running", containerName)

			case "restarting":
				log.Warnf("Container %s is restarting", containerName)

			case "paused":
				log.Warnf("Container %s is paused", containerName)

			case "dead":
				log.Errorf("Container %s is in dead state", containerName)

			default:
				log.Warnf("Container %s is in unknown state: %s", containerName, state)
			}
		}

		// Check for apps without running containers
		for appName := range appNameMap {
			status, exists := appStatusMap[appName]
			if !exists || status == "unknown" {
				log.Errorf("App %s has no running container", appName)

				// 查找应用记录
				appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.Name.Eq(appName)).First()
				if err != nil {
					log.Errorf("Failed to find app record for %s: %v", appName, err)
					continue
				}
				// 更新应用状态为错误
				_, err = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Not(repo.AppInstalled.Status.Eq(constant.UpErr)).Updates(model.AppInstalled{
					Status:  constant.Error,
					Message: "Container not found or in unknown state",
				})
				if err != nil {
					log.Errorf("Failed to update app status for %s: %v", appName, err)
				}

				// 记录日志
				// insertLog(appInstalled.ID, "容器监控", "容器不存在或状态未知")
				// 根据容器状态更新应用状态
				continue
			}
			switch status {
			case "running":
				// 容器正常运行，更新状态为 Running
				updateAppStatus(appName, constant.Running, "")
			case "exited":
				// 获取退出详情
				container, err := dm.client.ContainerInspect(dm.ctx, appName)
				if err != nil {
					log.Errorf("Failed to inspect container %s: %v", appName, err)
					continue
				}
				if container.State.ExitCode == 0 {
					// 正常停止
					updateAppStatus(appName, constant.Stopped, "Container stopped normally")
				} else {
					// 异常退出
					message := fmt.Sprintf("Container exited with code %d: %s",
						container.State.ExitCode, container.State.Error)
					updateAppStatus(appName, constant.Error, message)
				}
			case "restarting":
				updateAppStatus(appName, constant.Restarting, "Container is restarting")
			case "paused":
				updateAppStatus(appName, constant.Paused, "Container is paused")
			case "dead":
				updateAppStatus(appName, constant.Dead, "Container is in dead state")
			default:
				updateAppStatus(appName, constant.UnHealthy, fmt.Sprintf("Container is in unknown state: %s", status))
			}

		}

		return nil
	}

	// 添加定时任务并立即执行一次
	go func() {
		// 先执行一次任务
		if err := task(); err != nil {
			log.Printf("Error monitoring containers: %v", err)
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := task(); err != nil {
					log.Printf("Error monitoring containers: %v", err)
				}
			case <-dm.ctx.Done():
				return
			}
		}
	}()
}

// InitDockerMonitoring 初始化Docker监控
func InitDockerMonitoring() error {
	monitor, err := NewDockerMonitor()
	if err != nil {
		return fmt.Errorf(constant.ErrDockerMonitorInit, err)
	}

	// 设置监控间隔为1分钟
	monitor.StartMonitoring(1 * time.Minute)
	return nil
}

// 辅助函数：更新应用状态
func updateAppStatus(appName string, status string, message string) {
	appInstalled, err := repo.AppInstalled.Select(repo.AppInstalled.ID, repo.AppInstalled.Status).Where(repo.AppInstalled.Name.Eq(appName)).First()
	if err != nil {
		log.Errorf("Failed to find app record for %s: %v", appName, err)
		return
	}

	// 跳过处理 Installing 状态的应用
	if strings.EqualFold(appInstalled.Status, constant.Installing) {
		log.Debugf("Skipping status update for app %s as it is in Installing state", appName)
		return
	}

	// 只有状态发生变化时才更新
	if appInstalled.Status != status {
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
