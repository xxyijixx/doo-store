package req

type GenEnvReq struct {
	AppKey        string
	ContainerName string
	IPAddress     string
	Envs          map[string]any
	WriteFile     bool
}
