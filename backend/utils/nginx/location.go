package nginx

type LocationConfig struct {
	// 模板内容
	Template string
	// location名称, Key
	Name string
	// 代理服务器名称
	ProxyServerName string
	// 端口号
	Port int
	// 自定义配置项
	CustomOptions map[string]string
}

// NewLocationConfig 创建一个新的LocationConfig实例
func NewLocationConfig(name, proxyServer string) *LocationConfig {
	return &LocationConfig{
		Name:            name,
		ProxyServerName: proxyServer,
		CustomOptions:   make(map[string]string),
	}
}

// WithTemplate 设置模板
func (lc *LocationConfig) WithTemplate(tmpl string) *LocationConfig {
	lc.Template = tmpl
	return lc
}

// WithPort 设置端口
func (lc *LocationConfig) WithPort(port int) *LocationConfig {
	lc.Port = port
	return lc
}

// WithCustomOption 添加自定义配置项
func (lc *LocationConfig) WithCustomOption(key, value string) *LocationConfig {
	lc.CustomOptions[key] = value
	return lc
}