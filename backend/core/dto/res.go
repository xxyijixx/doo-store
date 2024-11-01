package dto

type PageResult struct {
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}

type Response struct {
	Code int         `json:"code" example:"200"`
	Msg  string      `json:"msg" example:"success"`
	Data interface{} `json:"data"`
}
