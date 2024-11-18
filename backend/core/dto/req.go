package dto

import "github.com/gin-gonic/gin"

type PageInfo struct {
	Page     int `form:"page" json:"page" binding:"required,number"`
	PageSize int `form:"page_size" json:"page_size" binding:"required,number"`
}

type ServiceContext struct {
	C        *gin.Context
	UserInfo UserInfoResp
}
