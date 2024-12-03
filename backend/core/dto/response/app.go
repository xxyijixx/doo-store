package response

import (
	"doo-store/backend/core/model"
)

type FormField struct {
	Default  string `json:"default,omitempty"`
	Label    string `json:"label"`
	EnvKey   string `json:"env_key,omitempty"`
	Key      string `json:"key,omitempty"`
	Value    string `json:"value,omitempty"`
	Values   any    `json:"values"`
	Type     string `json:"type"`
	Rule     string `json:"rule"`
	Required bool   `json:"required"`
}

type AppDetail struct {
	model.AppDetail
	Params AppParams `json:"params"`
}

type AppParams struct {
	FormFields []*FormField `json:"form_fields"`
}

type AppInstalledParamsResp struct {
	Params        []*FormField `json:"params"`
	DockerCompose string       `json:"docker_compose"`
	CPUS          string       `json:"cpus"`
	MemoryLimit   string       `json:"memory_limit"`
}

type GetInstalledPluginInfoResp struct {
	Name          string `json:"name"`
	Key           string `json:"key"`
	Location      string `json:"location"`
	Status        string `json:"status"`
	CloudProvider string `json:"cloud_provider,omitempty"`
}
