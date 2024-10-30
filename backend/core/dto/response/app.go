package response

type SyncRes struct {
	FormFields []FormField
}

type FormField struct {
	Default  string `json:"name"`
	Label    string `json:"label"`
	EnvKey   string `json:"envKey"`
	Type     string `json:"type"`
	Rule     string `json:"rule"`
	Edit     bool   `json:"edit"`
	Required bool   `json:"required"`
}
