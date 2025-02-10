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
	ErrPluginOnlyOneService          = "ErrPluginOnlyOneService"          // 只能有一个服务
	ErrPluginNotAllowedPrivileged    = "ErrPluginNotAllowedPrivileged"    // 不允许使用特权模式
	ErrPluginInvalidLocalVolumeMount = "ErrPluginInvalidLocalVolumeMount" // 本地卷挂载路径无效
	ErrPluginEnvVarInVolumeMount     = "ErrPluginEnvVarInVolumeMount"     // 不允许在挂载路径使用环境变量
	ErrPluginInstallFailed           = "ErrPluginInstallFailed"           // 插件安装失败
	ErrPluginUninstallFailed         = "ErrPluginUninstallFailed"         // 插件卸载失败
	ErrPluginParamInvalid            = "ErrPluginParamInvalid"            // 插件参数无效
	ErrPluginNotInstalled            = "ErrPluginNotInstalled"            // 插件未成功安装，请重新安装
	ErrPluginNotRunning              = "ErrPluginNotRunning"              // 插件未运行
	ErrPluginMissingParam            = "ErrPluginMissingParam"            // 缺少必填参数 {{.detail}}
	ErrPluginKeyExist                = "ErrPluginKeyExist"                // 插件key已存在
	ErrPluginUnsupportedAction       = "ErrPluginUnsupportedAction"       // 不支持的操作
	ErrPluginInfoFailed              = "ErrPluginInfoFailed"              // 获取插件信息失败
	ErrPluginVersionFailed           = "ErrPluginVersionFailed"           // 获取版本信息失败
	ErrPluginDependencyFailed        = "ErrPluginDependencyFailed"        // 检查依赖版本失败
	ErrPluginParamParseFailed        = "ErrPluginParamParseFailed"        // 解析插件参数失败
	ErrPluginModifyParamFailed       = "ErrPluginModifyParamFailed"       // 修改参数失败
	ErrPluginRestartFailed           = "ErrPluginRestartFailed"           // 插件重启失败

	// docker
	ErrDockerClientCreate     = "ErrDockerClientCreate"     // 创建Docker客户端失败
	ErrDockerListContainers   = "ErrDockerListContainers"   // 获取容器列表失败
	ErrDockerFindApps         = "ErrDockerFindApps"         // 查找应用失败
	ErrDockerMonitorInit      = "ErrDockerMonitorInit"      // 初始化Docker监控失败
	ErrDockerExecCreate       = "ErrDockerExecCreate"       // 创建执行命令失败
	ErrDockerExecAttach       = "ErrDockerExecAttach"       // 附加到执行命令失败
	ErrDockerExecRun          = "ErrDockerExecRun"          // 执行命令失败
	ErrNginxContainerNotFound = "ErrNginxContainerNotFound" // 未找到Nginx容器

	// nginx
	ErrNginxWriteFile    = "ErrNginxWriteFile"    // 写入文件失败
	ErrNginxParseContent = "ErrNginxParseContent" // 解析内容失败
	ErrNginxGetContainer = "ErrNginxGetContainer" // 获取Nginx容器失败

	// log
	ErrLogGetFailed  = "ErrLogGetFailed"  // 获取日志失败
	ErrLogReadFailed = "ErrLogReadFailed" // 读取日志失败

	// dootask
	ErrDooTaskDataFormat           = "ErrDooTaskDataFormat"           //数据格式错误
	ErrDooTaskResponseFormat       = "ErrDooTaskResponseFormat"       //响应格式错误
	ErrDooTaskRequestFailed        = "ErrDooTaskRequestFailed"        //请求失败
	ErrDooTaskUnmarshalResponse    = "ErrDooTaskUnmarshalResponse"    //解析响应失败：{{.detail}}
	ErrDooTaskRequestFailedWithErr = "ErrDooTaskRequestFailedWithErr" //请求失败：{{.detail}}
)
