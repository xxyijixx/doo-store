package task

import (
	"context"
	"doo-store/backend/core/repo"
	"doo-store/backend/utils/docker"
	"fmt"
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
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	return &DockerMonitor{
		client: cli.GetClient(),
		ctx:    context.Background(),
	}, nil
}

// StartMonitoring 开始监控任务
func (dm *DockerMonitor) StartMonitoring(interval time.Duration) {
	task := func() error {
		apps, err := repo.AppInstalled.Find()
		if err != nil {
			return fmt.Errorf("failed to find apps: %v", err)
		}

		// Create filter args with app names
		filterArgs := filters.NewArgs()
		appNameMap := make(map[string]bool)
		for _, app := range apps {
			filterArgs.Add("name", app.Name)
			appNameMap[app.Name] = true
		}

		// List all containers without filtering
		containers, err := dm.client.ContainerList(dm.ctx, container.ListOptions{
			All:     true,
			Filters: filterArgs,
		})
		if err != nil {
			return fmt.Errorf("failed to list containers: %v", err)
		}

		for _, container := range containers {
			// fmt.Printf("Container ID: %s\n", container.ID[:12])
			// fmt.Printf("Image: %s\n", container.Image)
			// fmt.Printf("Status: %s\n", container.Status)
			// fmt.Printf("State: %s\n", container.State)
			// fmt.Printf("Names: %v\n", container.Names)

			// Check if container name exists in apps
			containerName := container.Names[0][1:] // Remove leading slash
			if !appNameMap[containerName] {
				log.Error("Found container not in apps list: ", containerName)
				continue
			}

			state := container.State
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
		return fmt.Errorf("failed to initialize Docker monitor: %v", err)
	}

	// 设置监控间隔为1分钟
	monitor.StartMonitoring(1 * time.Minute)
	return nil
}
