package error

import (
	"doo-store/backend/i18n"
	"errors"

	"github.com/gin-gonic/gin"
)

type WithError struct {
	Msg     string
	Detail  interface{}
	Map     map[string]interface{}
	Err     error
	Content string
}

func (e WithError) Error() string {
	return e.Content
}

func (e WithError) GenContent(ctx *gin.Context) string {
	content := ""
	if e.Detail != nil {
		content = i18n.GetErrMsg(ctx, e.Msg, map[string]any{"detail": e.Detail})
	} else if e.Map != nil {
		content = i18n.GetErrMsg(ctx, e.Msg, e.Map)
	} else {
		content = i18n.GetErrMsg(ctx, e.Msg, nil)
	}
	if content == "" {
		if e.Err != nil {
			return e.Err.Error()
		}
		return errors.New(e.Msg).Error()
	}
	return content
}

func New(ctx *gin.Context, Key string) WithError {
	e := WithError{
		Msg:    Key,
		Detail: nil,
		Err:    nil,
	}
	e.Content = e.GenContent(ctx)
	return e
}

func WithDetail(ctx *gin.Context, Key string, detail any, err error) WithError {
	e := WithError{
		Msg:    Key,
		Detail: detail,
		Err:    err,
	}
	e.Content = e.GenContent(ctx)
	return e
}

func WithErr(ctx *gin.Context, Key string, err error) WithError {
	e := WithError{
		Msg:    Key,
		Detail: "",
		Err:    err,
	}
	e.Content = e.GenContent(ctx)
	return e
}

func WithMap(ctx *gin.Context, Key string, maps map[string]any, err error) WithError {
	e := WithError{
		Msg: Key,
		Map: maps,
		Err: err,
	}
	e.Content = e.GenContent(ctx)
	return e
}
