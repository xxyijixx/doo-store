package constant

var (
	ErrEnvProhibition   = "ErrEnvProhibition"   //当前环境禁止此操作
	ErrInvalidParameter = "ErrInvalidParameter" //参数错误
	ErrTypeNotLogin     = "ErrTypeNotLogin"     //未登录
	ErrRequestTimeout   = "ErrRequestTimeout"   //请求超时
	ErrNoPermission     = "ErrNoPermission"     //权限不足

	// plugin
	ErrPluginAdminNotCancel          = "ErrPluginAdminNotCancel"          // 仅限管理员操作
	ErrPluginVersionNotSupport       = "ErrPluginVersionNotSupport"       // 当前版本不满足要求，需要版本 {{.detail}} 或以上
	ErrPluginUnmarshalDockerCompose  = "ErrPluginUnmarshalDockerCompose"  // 无法解析 Docker Compose 文件
	ErrPluginNetworkModeHost         = "ErrPluginNetworkModeHost"         // 使用了 host 网络模式
	ErrPluginNotAllowedPrivileged    = "ErrPluginNotAllowedPrivileged"    // 不允许使用特权模式
	ErrPluginInvalidLocalVolumeMount = "ErrPluginInvalidLocalVolumeMount" // 本地卷挂载路径无效
	ErrPluginEnvVarInVolumeMount     = "ErrPluginEnvVarInVolumeMount"     // 不允许在挂载路径使用环境变量

	// dootask
	ErrDooTaskDataFormat           = "ErrDooTaskDataFormat"           //数据格式错误
	ErrDooTaskResponseFormat       = "ErrDooTaskResponseFormat"       //响应格式错误
	ErrDooTaskRequestFailed        = "ErrDooTaskRequestFailed"        //请求失败
	ErrDooTaskUnmarshalResponse    = "ErrDooTaskUnmarshalResponse"    //解析响应失败：{{.detail}}
	ErrDooTaskRequestFailedWithErr = "ErrDooTaskRequestFailedWithErr" //请求失败：{{.detail}}
)
