package helper

import (
	"doo-store/backend/core/dto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckBindAndValidate(req interface{}, c *gin.Context) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}
	return nil
}

func CheckBindQueryAndValidate(req interface{}, c *gin.Context) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return err
	}
	return nil
}

// Token 获取Token（Header、Query、Cookie）
func Token(c *gin.Context) string {
	token := c.GetHeader("token")
	if token == "" {
		token = Input(c, "token")
	}
	if token == "" {
		token = Cookie(c, "token")
	}
	return token
}

// Version 获取Version（Header、Query、Cookie）
func Version(c *gin.Context) string {
	token := c.GetHeader("version")
	if token == "" {
		token = Input(c, "version")
	}
	if token == "" {
		token = Cookie(c, "version")
	}
	return token
}

// Input 获取参数（优先POST、取Query）
func Input(c *gin.Context, key string) string {
	if c.PostForm(key) != "" {
		return strings.TrimSpace(c.PostForm(key))
	}
	return strings.TrimSpace(c.Query(key))
}

// Scheme 获取Scheme
func Scheme(c *gin.Context) string {
	scheme := "http://"
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https://"
	}
	return scheme
}

// Cookie 获取Cookie
func Cookie(c *gin.Context, name string) string {
	value, _ := c.Cookie(name)
	return value
}

func Response(c *gin.Context, code int, msg string, values ...any) {
	var data any
	if len(values) == 1 {
		data = values[0]
	} else if len(values) == 0 {
		data = gin.H{}
	} else {
		data = values
	}
	c.JSON(code, dto.Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	c.Abort()
}

func SuccessWith(ctx *gin.Context, values ...any) {
	Response(ctx, http.StatusOK, "success", values...)
}

// Error 失败
func Error(c *gin.Context, values ...any) {
	Response(c, http.StatusBadRequest, "error", values...)
}

func ErrorWith(c *gin.Context, msgKey string, err error, values ...any) {
	Response(c, http.StatusBadRequest, "error", values...)
}
