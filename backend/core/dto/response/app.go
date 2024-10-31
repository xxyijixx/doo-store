package response

type FormField struct {
	Default  string `json:"default"`
	Label    string `json:"label"`
	EnvKey   string `json:"envKey"`
	Type     string `json:"type"`
	Rule     string `json:"rule"`
	Required bool   `json:"required"`
}
