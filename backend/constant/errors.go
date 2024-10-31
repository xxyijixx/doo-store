package constant

var (
	ErrNoPermission = "ErrNoPermission"

	// dootask
	ErrDooTaskDataFormat           = "ErrDooTaskDataFormat"           //数据格式错误
	ErrDooTaskResponseFormat       = "ErrDooTaskResponseFormat"       //响应格式错误
	ErrDooTaskRequestFailed        = "ErrDooTaskRequestFailed"        //请求失败
	ErrDooTaskUnmarshalResponse    = "ErrDooTaskUnmarshalResponse"    //解析响应失败：{{.detail}}
	ErrDooTaskRequestFailedWithErr = "ErrDooTaskRequestFailedWithErr" //请求失败：{{.detail}}
)
