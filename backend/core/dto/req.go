package dto

type PageInfo struct {
	Page     int `form:"page" json:"page" validate:"required,number"`
	PageSize int `form:"pageSize" json:"pageSize" validate:"required,number"`
}
