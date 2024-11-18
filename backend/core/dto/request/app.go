package request

import "doo-store/backend/core/dto"

type AppSearch struct {
	dto.PageInfo
	ID          int64  `form:"id" json:"id"`
	Name        string `form:"name" json:"name"`
	Class       string `form:"class" json:"class"`
	Description string `form:"description" json:"description"`
}

type AppInstall struct {
	Name          string                 `json:"-"`
	Key           string                 `json:"-"`
	InstalledId   int64                  `json:"-"`
	DockerCompose string                 `json:"docker_compose" binding:"required"`
	CPUS          string                 `json:"cpus" binding:"required"`
	MemoryLimit   string                 `json:"memory_limit" binding:"required"`
	Params        map[string]interface{} `json:"params" binding:"required"`
}

type AppUnInstall struct {
	Key string `json:"-"`
}

type AppInstalledOperate struct {
	Action string                 `json:"action"`
	Key    string                 `json:"-"`
	Params map[string]interface{} `json:"params"`
}

type AppInstalledSearch struct {
	dto.PageInfo
	Name        string `form:"name" json:"name"`
	Class       string `form:"class" json:"class"`
	Description string `form:"description" json:"description"`
}

type AppLogsSearch struct {
	Id    int64  `json:"-"`
	Since string `form:"since"`
	Until string `form:"until"`
	Tail  int    `form:"tail"`
}
