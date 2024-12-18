package dto

import (
	"doo-store/backend/config"
	"doo-store/backend/core/dto/response"
	"encoding/json"
	"fmt"
	"strings"
)

type Plugin struct {
	Name           string       `json:"name"`
	Key            string       `json:"key"`
	Description    string       `json:"description"`
	Icon           string       `json:"icon"`
	Version        string       `json:"version"`
	Github         string       `json:"github"`
	Class          string       `json:"class"`
	DependsVersion string       `json:"depends_version"`
	Repo           string       `json:"repo"`
	Volume         []Volume     `json:"volume"`
	Env            []EnvElement `json:"env"`
	Command        string       `json:"command"`
	NginxConfig    string       `json:"nginx_config"`
	DockerCompose  string       `json:"docker_compose"`
}

type EnvElement struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type Volume struct {
	Local  string `json:"local"`
	Target string `json:"target"`
}

func (p *Plugin) GenComposeFile() string {
	if p.DockerCompose != "" {
		return p.DockerCompose
	}
	composeContent := p.GenService()
	composeContent += p.GenNetwork()
	return composeContent
}

func (p *Plugin) GenNetwork() string {
	networkContent := make([]string, 0)
	networkContent = append(networkContent, "networks:")
	networkContent = append(networkContent, fmt.Sprintf("%s%s:", p.getSpaces(1), config.EnvConfig.GetNetworkName()))
	networkContent = append(networkContent, fmt.Sprintf("%sexternal: true", p.getSpaces(2)))

	networkContent = append(networkContent, "\n")

	return strings.Join(networkContent, "\n")
}

func (p *Plugin) GenService() string {
	serviceContent := make([]string, 0)
	serviceContent = append(serviceContent, "services:")
	serviceContent = append(serviceContent, fmt.Sprintf("%s%s:", p.getSpaces(1), p.Key))
	serviceContent = append(serviceContent, fmt.Sprintf("%simage: %s:%s", p.getSpaces(2), p.Repo, p.Version))
	serviceContent = append(serviceContent, fmt.Sprintf("%srestart: always", p.getSpaces(2)))
	serviceContent = append(serviceContent, fmt.Sprintf("%scontainer_name: ${CONTAINER_NAME}", p.getSpaces(2)))
	// networks:
	serviceContent = append(serviceContent, fmt.Sprintf("%snetworks:", p.getSpaces(2)))
	serviceContent = append(serviceContent, fmt.Sprintf("%s%s:", p.getSpaces(3), config.EnvConfig.GetNetworkName()))
	serviceContent = append(serviceContent, fmt.Sprintf("%sipv4_address: ${IP_ADDRESS}", p.getSpaces(4)))

	if len(p.Volume) != 0 {
		serviceContent = append(serviceContent, fmt.Sprintf("%svolumes:", p.getSpaces(2)))
		for _, v := range p.Volume {
			serviceContent = append(serviceContent, fmt.Sprintf("%s- %s:%s", p.getSpaces(3), v.Local, v.Target))
		}
	}
	if len(p.Env) != 0 {
		serviceContent = append(serviceContent, fmt.Sprintf("%senvironment:", p.getSpaces(2)))
		for _, env := range p.Env {
			serviceContent = append(serviceContent, fmt.Sprintf("%s- %s=${%s}", p.getSpaces(3), env.Key, env.Key))
		}
	}

	serviceContent = append(serviceContent, fmt.Sprintf("%scpus: \"${CPUS}\"", p.getSpaces(2)))
	serviceContent = append(serviceContent, fmt.Sprintf("%smem_limit: \"${MEMORY_LIMIT}\"", p.getSpaces(2)))

	serviceContent = append(serviceContent, fmt.Sprintf("%slabels:", p.getSpaces(2)))
	serviceContent = append(serviceContent, fmt.Sprintf("%screatedBy: \"Apps\"", p.getSpaces(3)))

	if p.Command != "" {
		serviceContent = append(serviceContent, fmt.Sprintf("%scommand: %s", p.getSpaces(2), p.Command))
	}

	serviceContent = append(serviceContent, "\n")

	return strings.Join(serviceContent, "\n")
}

func (p *Plugin) GenParams() string {
	formFields := make([]response.FormField, 0)
	for _, env := range p.Env {
		formField := response.FormField{
			Label:    env.Name,
			Default:  fmt.Sprintf("%v", env.Value),
			EnvKey:   env.Key,
			Type:     env.Type,
			Required: env.Required,
		}
		formFields = append(formFields, formField)
	}
	params := map[string]interface{}{
		"form_fields": formFields,
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// GenNginxConfig 生成Nginx配置信息
func (p *Plugin) GenNginxConfig() string {
	if p.NginxConfig != "" {
		return p.NginxConfig
	}
	return `location /plugin/{{.Key}}/ {
	proxy_http_version 1.1;
	proxy_set_header X-Real-IP $remote_addr;
	proxy_set_header X-Real-PORT $remote_port;
	proxy_set_header X-Forwarded-Host $the_host;
	proxy_set_header X-Forwarded-Proto $the_scheme;
	proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	proxy_set_header Host $http_host;
	proxy_set_header Scheme $scheme;
	proxy_set_header Server-Protocol $server_protocol;
	proxy_set_header Server-Name $server_name;
	proxy_set_header Server-Addr $server_addr;
	proxy_set_header Server-Port $server_port;
	proxy_set_header Upgrade $http_upgrade;
	proxy_set_header Connection $connection_upgrade;
	proxy_pass http://{{.ContainerName}}:{{.Port}}/;
}`
}

func (Plugin) getSpaces(num int) string {
	spaces := strings.Repeat(" ", 2*num)
	return spaces
}
