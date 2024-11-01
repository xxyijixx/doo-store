package dto

type PageInfo struct {
	Page     int `form:"page" json:"page" validate:"required,number"`
	PageSize int `form:"page_size" json:"page_size" validate:"required,number"`
}
