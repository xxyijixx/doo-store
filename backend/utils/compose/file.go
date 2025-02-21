package compose

import (
	"doo-store/backend/constant"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type DockerComposeConfig struct {
	Version  string                   `yaml:"version"`
	Services map[string]ServiceConfig `yaml:"services"`
	Networks map[string]NetworkConfig `yaml:"networks,omitempty"`
	Volumes  map[string]VolumeConfig  `yaml:"volumes,omitempty"`
	Configs  map[string]ConfigConfig  `yaml:"configs,omitempty"`
	Secrets  map[string]SecretConfig  `yaml:"secrets,omitempty"`
}

type ServiceConfig struct {
	Image         string                     `yaml:"image,omitempty"`
	Restart       string                     `yaml:"restart,omitempty"`
	ContainerName string                     `yaml:"container_name,omitempty"`
	Ports         StringOrList               `yaml:"ports,omitempty"`
	Env           MapOrSlice                 `yaml:"environment,omitempty"`
	Volumes       StringOrList               `yaml:"volumes,omitempty"`
	NetworkMode   string                     `yaml:"network_mode,omitempty"`
	Networks      map[string]NetworkSettings `yaml:"networks,omitempty"`
	Privileged    bool                       `yaml:"privileged,omitempty"`
	Command       StringOrList               `yaml:"command,omitempty"`
	Entrypoint    StringOrList               `yaml:"entrypoint,omitempty"`
	DependsOn     StringOrList               `yaml:"depends_on,omitempty"`
	ExtraHosts    StringOrList               `yaml:"extra_hosts,omitempty"`
	DNS           StringOrList               `yaml:"dns,omitempty"`
	DNSSearch     StringOrList               `yaml:"dns_search,omitempty"`
	Labels        map[string]string          `yaml:"labels,omitempty"`
	Logging       LoggingConfig              `yaml:"logging,omitempty"`
	HealthCheck   HealthCheckConfig          `yaml:"healthcheck,omitempty"`
	Deploy        DeployConfig               `yaml:"deploy,omitempty"`
	CapAdd        StringOrList               `yaml:"cap_add,omitempty"`
	CapDrop       StringOrList               `yaml:"cap_drop,omitempty"`
	WorkingDir    string                     `yaml:"working_dir,omitempty"`
	User          string                     `yaml:"user,omitempty"`
	Sysctls       map[string]string          `yaml:"sysctls,omitempty"`
	Ulimits       UlimitsConfig              `yaml:"ulimits,omitempty"`
	Build         BuildConfig                `yaml:"build,omitempty"`
}

type NetworkConfig struct {
	External bool `yaml:"external"`
}

type NetworkSettings struct {
	IPAddress string `yaml:"ipv4_address,omitempty"` // 添加静态IP地址
}

type VolumeConfig struct {
	Driver     string            `yaml:"driver,omitempty"`
	DriverOpts map[string]string `yaml:"driver_opts,omitempty"`
	External   bool              `yaml:"external,omitempty"`
}

type ConfigConfig struct {
	File string `yaml:"file"`
}

type SecretConfig struct {
	File string `yaml:"file"`
}

type LoggingConfig struct {
	Driver  string     `yaml:"driver"`
	Options MapOrSlice `yaml:"options,omitempty"`
}

type HealthCheckConfig struct {
	Test        StringOrList `yaml:"test,omitempty"`
	Interval    string       `yaml:"interval,omitempty"`
	Timeout     string       `yaml:"timeout,omitempty"`
	Retries     int          `yaml:"retries,omitempty"`
	StartPeriod string       `yaml:"start_period,omitempty"`
}

type DeployConfig struct {
	Replicas      int                 `yaml:"replicas,omitempty"`
	Resources     ResourcesConfig     `yaml:"resources,omitempty"`
	RestartPolicy RestartPolicyConfig `yaml:"restart_policy,omitempty"`
}

type ResourcesConfig struct {
	Limits       ResourceLimits `yaml:"limits,omitempty"`
	Reservations ResourceLimits `yaml:"reservations,omitempty"`
}

type ResourceLimits struct {
	CPUs   string `yaml:"cpus,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

type RestartPolicyConfig struct {
	Condition   string `yaml:"condition,omitempty"`
	Delay       string `yaml:"delay,omitempty"`
	MaxAttempts int    `yaml:"max_attempts,omitempty"`
	Window      string `yaml:"window,omitempty"`
}

type UlimitsConfig struct {
	Nofile UlimitValues `yaml:"nofile,omitempty"`
}

type UlimitValues struct {
	Soft int `yaml:"soft"`
	Hard int `yaml:"hard"`
}

type BuildConfig struct {
	Context    string       `yaml:"context,omitempty"`
	Dockerfile string       `yaml:"dockerfile,omitempty"`
	Args       MapOrSlice   `yaml:"args,omitempty"`
	CacheFrom  StringOrList `yaml:"cache_from,omitempty"`
	Target     string       `yaml:"target,omitempty"`
}

// 自定义类型用于处理两种格式,处理 map[string]strin] 或 []string 兼容解析
type MapOrSlice map[string]string

// UnmarshalYAML 方法处理 []string 和 map[string]string 格式
func (e *MapOrSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*e = make(map[string]string)

	// 解析 map[string]string
	var mapFormat map[string]string
	if err := unmarshal(&mapFormat); err == nil {
		*e = mapFormat
		return nil
	}

	// 解析 []string
	var sliceFormat []string
	if err := unmarshal(&sliceFormat); err == nil {
		for _, item := range sliceFormat {
			parts := strings.SplitN(item, "=", 2)
			key := parts[0]
			value := ""
			if len(parts) == 2 {
				value = parts[1]
			}
			(*e)[key] = value
		}
		return nil
	}

	return fmt.Errorf("unsupported environment format")
}

func (e MapOrSlice) MarshalYAML() (interface{}, error) {
	var envList []string
	for key, value := range e {
		if value == "" {
			envList = append(envList, key)
		} else {
			envList = append(envList, fmt.Sprintf("%s=%s", key, value))
		}
	}
	return envList, nil
}

// StringOrList 处理 string 或 []string 兼容解析
type StringOrList []string

func (s *StringOrList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var list []string
	if err := unmarshal(&list); err == nil {
		*s = list
		return nil
	}

	var single string
	if err := unmarshal(&single); err == nil {
		*s = []string{single}
		return nil
	}

	return fmt.Errorf("unsupported format for StringOrList")
}

func (s StringOrList) MarshalYAML() (interface{}, error) {
	if len(s) == 1 {
		return s[0], nil
	}
	return []string(s), nil
}

func PreCheck(content string) (*DockerComposeConfig, error) {
	var config DockerComposeConfig
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return nil, errors.New(constant.ErrPluginUnmarshalDockerCompose)
	}

	return &config, config.preCheck()
}

func FullCheck(content string, envContent string) (*DockerComposeConfig, error) {
	envMap, err := parseEnvContent(envContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse envContent: %w", err)
	}

	// 替换 content 中的环境变量
	result := replaceEnvVars(content, envMap)
	var config DockerComposeConfig
	err = yaml.Unmarshal([]byte(result), &config)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return nil, errors.New(constant.ErrPluginUnmarshalDockerCompose)
	}
	return &config, config.fullCheck()
}

// DockerCompose 文件检查
func (dcc DockerComposeConfig) preCheck() error {
	for _, serviceConfig := range dcc.Services {
		if serviceConfig.Privileged {
			return errors.New(constant.ErrPluginNotAllowedPrivileged)
		}
		if serviceConfig.NetworkMode == "host" {
			return errors.New(constant.ErrPluginNetworkModeHost)
		}
		// 当前只允许一个service
		// if len(dcc.Services) > 1 {
		// 	return errors.New(constant.ErrPluginOnlyOneService)
		// }
	}
	return nil
}

// 对最终执行的 Docker Compose 文件进行全量检查
func (dcc DockerComposeConfig) fullCheck() error {
	dcc.preCheck() // 先进行预检查

	return nil
}

// 检测 Docker Compose 文件中的挂载目录
// func checkLocalVolumeMounts(volumes []string) error {
// 	// 定义正则表达式，用来匹配以 ./ 开头的挂载路径，避免包含 `..`
// 	re := regexp.MustCompile(`^\.\/[^:]+$`)

// 	for _, volume := range volumes {
// 		// 拆分挂载的路径和目标路径
// 		parts := strings.Split(volume, ":")
// 		if len(parts) > 1 {
// 			localPath := parts[0] // 本地路径

// 			// 检查本地路径是否以 ./ 开头，并且不包含返回上级目录的 `..`
// 			if !re.MatchString(localPath) {
// 				return errors.New(constant.ErrPluginInvalidLocalVolumeMount)
// 			}

// 			// 检查是否有环境变量
// 			if strings.Contains(localPath, "${") {
// 				return errors.New(constant.ErrPluginEnvVarInVolumeMount)
// 			}
// 		}
// 	}
// 	return nil
// }

// 提取 Docker Compose 文件中的 IP 地址
func (dcc *DockerComposeConfig) ExtractIpAddress() []string {
	var ipList []string
	for _, serviceConfig := range dcc.Services {
		for _, networkConfig := range serviceConfig.Networks {
			ipList = append(ipList, networkConfig.IPAddress)
		}
	}
	return ipList
}

// 提取 Docker Compose 文件中的 IP 地址
func (dcc *DockerComposeConfig) ExtractContainerName() []string {
	var containerNameList []string
	for _, serviceConfig := range dcc.Services {
		containerNameList = append(containerNameList, serviceConfig.ContainerName)
	}
	return containerNameList
}
