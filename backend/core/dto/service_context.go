package dto

import (
	"github.com/gin-gonic/gin"
)

// ServiceContext 服务上下文，用于传递请求上下文信息
type ServiceContext struct {
	C           *gin.Context  // Gin 上下文
	UserInfo    *UserInfoResp // 用户信息
	RequestID   string        // 请求ID，用于追踪
	Language    string        // 国际化语言
	ClientIP    string        // 客户端IP
	UserAgent   string        // 用户代理
	RequestTime int64         // 请求时间戳
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c *gin.Context) ServiceContext {
	// 获取用户信息
	var userInfo *UserInfoResp
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*UserInfoResp); ok {
			userInfo = u
		}
	}

	// 获取语言设置
	lang := c.GetHeader("language")
	if lang == "" {
		lang = "zh" // 默认中文
	}

	return ServiceContext{
		C:           c,
		UserInfo:    userInfo,
		RequestID:   c.GetString("X-Request-ID"),
		Language:    lang,
		ClientIP:    c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
		RequestTime: c.GetInt64("request_time"),
	}
}

// GetUserID 获取用户ID
func (ctx *ServiceContext) GetUserID() int {
	if ctx.UserInfo != nil && ctx.UserInfo.UserBasicResp != nil {
		return ctx.UserInfo.Userid
	}
	return 0
}

// IsAdmin 判断当前用户是否是管理员
func (ctx *ServiceContext) IsAdmin() bool {
	if ctx.UserInfo != nil {
		return ctx.UserInfo.IsAdmin()
	}
	return false
}

// GetLanguage 获取当前语言设置
func (ctx *ServiceContext) GetLanguage() string {
	return ctx.Language
}

// GetClientIP 获取客户端IP
func (ctx *ServiceContext) GetClientIP() string {
	return ctx.ClientIP
}

// GetUserAgent 获取用户代理
func (ctx *ServiceContext) GetUserAgent() string {
	return ctx.UserAgent
}

// GetRequestID 获取请求ID
func (ctx *ServiceContext) GetRequestID() string {
	return ctx.RequestID
}

// GetRequestTime 获取请求时间
func (ctx *ServiceContext) GetRequestTime() int64 {
	return ctx.RequestTime
}

// GetContext 获取原始的 gin.Context
func (ctx *ServiceContext) GetContext() *gin.Context {
	return ctx.C
}
