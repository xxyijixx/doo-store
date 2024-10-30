package helper

import (
	"doo-store/backend/constant"
	"doo-store/backend/core/dto"
	"net/http"

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

func SuccessWithData(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := dto.Response{
		Code: constant.CodeSuccess,
		Data: data,
	}
	ctx.JSON(http.StatusOK, res)
	ctx.Abort()
}
