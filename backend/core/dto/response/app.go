package response

import "doo-store/backend/core/model"

type FormField struct {
	Default  string `json:"default"`
	Label    string `json:"label"`
	EnvKey   string `json:"env_key"`
	Type     string `json:"type"`
	Rule     string `json:"rule"`
	Required bool   `json:"required"`
}

type AppDetail struct {
	model.AppDetail
	Params AppParams `json:"params"`
}

type AppParams struct {
	FormFields []FormField `json:"form_fields"`
}
