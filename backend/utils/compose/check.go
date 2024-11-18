package compose

import (
	"doo-store/backend/constant"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// 自定义类型用于处理两种格式
type Environment map[string]string

// UnmarshalYAML 方法处理 []string 和 map[string]string 格式
func (e *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// 尝试解析为 map[string]string 格式
	var mapFormat map[string]string
	if err := unmarshal(&mapFormat); err == nil {
		*e = mapFormat
		return nil
	}

	// 尝试解析为 []string 格式
	var sliceFormat []string
	if err := unmarshal(&sliceFormat); err == nil {
		*e = make(map[string]string)
		for _, item := range sliceFormat {
			parts := strings.SplitN(item, "=", 2)
			if len(parts) == 2 {
				(*e)[parts[0]] = parts[1]
			}
		}
		return nil
	}

	return fmt.Errorf("environment format not supported")
}

type DockerComposeConfig struct {
	Version  string                   `yaml:"version"`
	Services map[string]ServiceConfig `yaml:"services"`
	Networks map[string]NetworkConfig `yaml:"networks,omitempty"`
}

type ServiceConfig struct {
	Image         string      `yaml:"image"`
	Restart       string      `yaml:"restart,omitempty"`
	ContainerName string      `yaml:"container_name"`
	Ports         []string    `yaml:"ports,omitempty"`
	Env           Environment `yaml:"environment,omitempty"`
	Volumes       []string    `yaml:"volumes,omitempty"`
	NetworkMode   string      `yaml:"network_mode"`
	Privileged    bool        `yaml:"privileged,omitempty"`
}

type NetworkConfig struct {
	External bool `yaml:"external"`
}

func Check(content string) error {
	var config DockerComposeConfig
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return errors.New(constant.ErrPluginUnmarshalDockerCompose)
	}
	for _, serviceConfig := range config.Services {
		if serviceConfig.Privileged {
			return errors.New(constant.ErrPluginNotAllowedPrivileged)
		}
		if serviceConfig.NetworkMode == "host" {
			return errors.New(constant.ErrPluginNetworkModeHost)
		}
		if err = checkLocalVolumeMounts(serviceConfig.Volumes); err != nil {
			return err
		}
	}
	return nil
}

// 检测 Docker Compose 文件中的挂载目录
func checkLocalVolumeMounts(volumes []string) error {
	// 定义正则表达式，用来匹配以 ./ 开头的挂载路径，避免包含 `..`
	re := regexp.MustCompile(`^\.\/[^:]+$`)

	for _, volume := range volumes {
		// 拆分挂载的路径和目标路径
		parts := strings.Split(volume, ":")
		if len(parts) > 1 {
			localPath := parts[0] // 本地路径

			// 检查本地路径是否以 ./ 开头，并且不包含返回上级目录的 `..`
			if !re.MatchString(localPath) {
				return errors.New(constant.ErrPluginInvalidLocalVolumeMount)
			}

			// 检查是否有环境变量
			if strings.Contains(localPath, "${") {
				return errors.New(constant.ErrPluginEnvVarInVolumeMount)
			}
		}
	}
	return nil
}
