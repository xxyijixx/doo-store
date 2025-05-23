package service

import (
	"context"
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/core/dto"
	"doo-store/backend/core/dto/request"
	"doo-store/backend/core/dto/response"
	"doo-store/backend/core/model"
	"doo-store/backend/core/repo"
	"doo-store/backend/utils/common"
	"doo-store/backend/utils/compose"
	"doo-store/backend/utils/docker"
	e "doo-store/backend/utils/error"
	"doo-store/backend/utils/nginx"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"

	schemasReq "doo-store/backend/core/schemas/req"

	"github.com/docker/docker/api/types/container"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AppInstallProcess struct {
	ctx                  dto.ServiceContext
	app                  *model.App
	appDetail            *model.AppDetail
	appInstalled         *model.AppInstalled
	appKey               string
	containerName        string
	defaultContainerName string
	envContent           string
	envJson              string
	req                  request.AppInstall
	ipAddress            string
	client               docker.Client
	dockerCompose        *compose.DockerComposeConfig
	finalDockerCompose   *compose.DockerComposeConfig
	nm                   *nginx.NginxManager
}

// NewAppInstallProcess 创建新的应用安装流程实例
func NewAppInstallProcess(ctx dto.ServiceContext, req request.AppInstall) *AppInstallProcess {
	return &AppInstallProcess{
		ctx: ctx,
		req: req,
	}
}

// ValidateInstallRequirements 验证安装要求
// 检查应用信息、版本依赖、是否已安装等
func (p *AppInstallProcess) ValidateInstallRequirements() error {
	var err error
	log.Info("开始验证应用安装要求")
	p.app, err = repo.App.Where(repo.App.Key.Eq(p.req.Key)).First()
	if err != nil {
		log.Error("查询应用信息失败:", err)
		return errors.New(constant.ErrPluginInfoFailed)
	}

	// 检测版本
	dootaskService := NewIDootaskService()
	versionInfoResp, err := dootaskService.GetVersoinInfo()
	if err != nil {
		log.Error("获取版本信息失败:", err)
		return errors.New(constant.ErrPluginVersionFailed)
	}

	check, err := versionInfoResp.CheckVersion(p.app.DependsVersion)
	if err != nil {
		log.Error("检测版本失败:", err)
		return errors.New(constant.ErrPluginDependencyFailed)
	}

	// 依赖版本不符合要求
	if !check {
		log.Warn("版本依赖不满足要求:", p.app.DependsVersion)
		return e.NewErrorWithMap(p.ctx.C, constant.ErrPluginVersionNotSupport, map[string]interface{}{
			"detail": p.app.DependsVersion,
		}, nil)
	}

	// 判断是否已安装
	log.Info("检查应用是否已安装")
	p.appInstalled, err = repo.AppInstalled.
		Select(repo.AppInstalled.ID, repo.AppInstalled.AppID).
		Where(repo.AppInstalled.AppID.Eq(p.app.ID)).
		First()

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error("查询已安装应用失败:", err)
			return errors.New(constant.ErrPluginInstallFailed)
		}
	}
	if p.appInstalled != nil {
		log.Warn("应用已安装")
		return errors.New(constant.ErrPluginInstallFailed)
	}

	log.Info("查询应用详细信息")
	p.appDetail, err = repo.AppDetail.Select(
		repo.AppDetail.ID,
		repo.AppDetail.AppID,
		repo.AppDetail.Repo,
		repo.AppDetail.Version,
		repo.AppDetail.Params,
		repo.AppDetail.DependsVersion,
		repo.AppDetail.NginxConfig,
	).Where(repo.AppDetail.AppID.Eq(p.app.ID)).First()
	if err != nil {
		log.Error("查询应用详细信息失败:", err)
		return errors.New(constant.ErrPluginInfoFailed)
	}
	log.Info("验证安装要求完成")
	return nil
}

// 分配IP地址
// 获取所有容器使用的IP并分配新IP
func (p *AppInstallProcess) DHCP() error {
	var err error
	p.client, err = docker.NewClient()
	if err != nil {
		log.Error("创建Docker客户端失败:", err)
		return err
	}
	client := p.client.GetClient()
	containers, err := client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Error("获取容器列表失败:", err)
		return fmt.Errorf(constant.ErrDockerListContainers, err)
	}

	// 检查所有容器使用的IP
	for _, container := range containers {
		if container.NetworkSettings != nil {
			for _, network := range container.NetworkSettings.Networks {
				if network.IPAddress != "" {
					if err := docker.GlobalIPAllocator.RegisterIP(network.IPAddress); err != nil {
						log.Debugf("注册IP失败 %s: %v", network.IPAddress, err)
					}
				}
			}
		}
	}

	// 分配新IP
	p.ipAddress, err = docker.GlobalIPAllocator.AllocateIP()
	if err != nil {
		log.Error("分配IP地址失败:", err)
		return err
	}
	log.Info("分配IP流程完成, 分配的IP:", p.ipAddress)
	return nil
}

func (p *AppInstallProcess) genEnv() error {
	var err error
	// 资源限制
	p.req.Params[constant.CPUS] = p.req.CPUS
	p.req.Params[constant.MemoryLimit] = p.req.MemoryLimit

	p.envContent, p.envJson, err = pluginHelper.GenEnv(schemasReq.GenEnvReq{
		AppKey:        p.appKey,
		ContainerName: p.defaultContainerName,
		IPAddress:     p.ipAddress,
		Envs:          p.req.Params,
		WriteFile:     false,
	})
	if err != nil {
		log.Error("生成环境变量失败:", err)
		return err
	}

	// 替换环境变量，再执行一次检查
	p.finalDockerCompose, err = compose.FullCheck(p.req.DockerCompose, p.envContent)
	if err != nil {
		return err
	}
	envChange := false
	// 是否释放原分配的IP并注册新IP
	ipList := p.finalDockerCompose.ExtractIpAddress()
	if len(ipList) > 0 {
		if ipList[0] != p.ipAddress {
			// 释放已注册的IP
			docker.GlobalIPAllocator.ReleaseIP(p.ipAddress)
			p.ipAddress = ipList[0]
			docker.GlobalIPAllocator.RegisterIP(p.ipAddress)
			envChange = true
		}
	}

	containerNameList := p.finalDockerCompose.ExtractContainerName()
	if len(containerNameList) > 0 {
		if containerNameList[0] != p.containerName {
			p.containerName = containerNameList[0]
		}
	}

	// 重新生成一下环境变量配置
	if envChange {
		p.envContent, p.envJson, err = pluginHelper.GenEnv(schemasReq.GenEnvReq{
			AppKey:        p.appKey,
			ContainerName: p.defaultContainerName,
			IPAddress:     p.ipAddress,
			Envs:          p.req.Params,
			WriteFile:     false,
		})
		if err != nil {
			log.Error("生成环境变量失败:", err)
			return err
		}
	}
	return nil
}

// ValidateParam 验证参数
// 检查docker-compose文件、创建工作目录、验证必填参数等
func (p *AppInstallProcess) ValidateParam() error {
	var err error

	// 检测 docker-compose 文件
	p.dockerCompose, err = compose.PreCheck(p.req.DockerCompose)
	if err != nil {
		log.Error("docker-compose文件验证失败:", err)
		return err
	}

	p.appKey = pluginHelper.GetAppKey(p.app.Key)

	// 创建工作目录
	workspaceDir := path.Join(constant.AppInstallDir, p.appKey)
	log.Info("创建工作目录:", workspaceDir)
	err = common.CreateDir(workspaceDir)
	if err != nil {
		log.Error("创建工作目录失败:", err)
		return err
	}

	// 容器名称
	p.defaultContainerName = config.EnvConfig.GetDefaultContainerName(p.app.Key)
	p.containerName = p.defaultContainerName

	paramJson, err := json.Marshal(p.req.Params)
	if err != nil {
		log.Error("参数序列化失败:", err)
		return err
	}

	params := response.AppParams{}
	err = common.StrToStruct(p.appDetail.Params, &params)
	if err != nil {
		log.Error("解析参数失败:", err)
		return errors.New(constant.ErrPluginParamParseFailed)
	}

	// 验证必填参数
	// for _, param := range params.FormFields {
	// 	if param.Required {
	// 		if _, exists := p.req.Params[param.EnvKey]; !exists {
	// 			log.Warn("缺少必填参数:", param.EnvKey)
	// 			return e.NewErrorWithDetail(p.ctx.C, constant.ErrPluginMissingParam, param.EnvKey, nil)
	// 		}
	// 	}
	// }
	vErr := dto.ValidateFormData(params.FormFields, p.req.Params)
	if len(vErr) > 0 {
		log.Warn("参数验证失败:", vErr)
		return vErr[0]
	}

	if err = p.genEnv(); err != nil {
		return err
	}

	p.appInstalled = &model.AppInstalled{
		Name:          p.containerName,
		AppID:         p.app.ID,
		AppDetailID:   p.appDetail.ID,
		Class:         p.app.Class,
		Repo:          p.appDetail.Repo,
		Version:       p.appDetail.Version,
		Params:        string(paramJson),
		Env:           p.envJson,
		DockerCompose: p.req.DockerCompose,
		Key:           p.app.Key,
		Status:        model.PluginStatusInstalling,
		IpAddress:     p.ipAddress,
	}
	// 更新插件状态
	err = repo.DB.Transaction(func(tx *gorm.DB) error {
		_, err = repo.Use(tx).App.Where(repo.App.ID.Eq(p.appInstalled.AppID)).Update(repo.App.Status, model.AppInUse)
		if err != nil {
			return err
		}
		err = repo.Use(tx).AppInstalled.Create(p.appInstalled)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error("更新应用状态失败:", err)
		return err
	}
	log.Info("参数验证完成")
	return nil
}

// Install 执行安装
func (p *AppInstallProcess) Install() error {
	log.Info("开始安装应用:", p.app.Name)
	var err error
	if p.appInstalled == nil {
		log.Error("未找到安装信息")
		return errors.New(constant.ErrPluginInstallFailed)
	}
	// TODO 安装服务
	appServiceList := make([]*model.AppServiceStatus, 0)
	for name, service := range p.finalDockerCompose.Services {
		IPAddress := []string{}
		for _, network := range service.Networks {
			IPAddress = append(IPAddress, network.IPAddress)
		}
		appService := model.AppServiceStatus{
			ServiceName:   name,
			ContainerName: service.ContainerName,
			IpAddress:     strings.Join(IPAddress, ","),
			Image:         service.Image,
			InstallID:     p.appInstalled.ID,
			Status:        model.PluginStatusInstalling,
		}
		appServiceList = append(appServiceList, &appService)
	}
	repo.AppServiceStatus.Create(appServiceList...)

	err = pluginActionManager.Up(p.appInstalled, p.envContent)
	if err != nil {
		log.Error("应用启动失败:", err)
		return err
	}
	log.Info("应用安装完成")
	return nil
}

// AddNginx 添加Nginx配置
// 插件安装的时候，需要向Nginx添加一个配置，如果添加配置失败，会将插件停止
func (p *AppInstallProcess) AddNginx() error {
	log.Info("开始配置Nginx")
	var err error
	p.nm, err = nginx.NewNginxManager()
	if err != nil {
		log.Error("创建Nginx管理器失败:", err)
		return err
	}

	port, err := p.client.GetImageFirstExposedPortByName(fmt.Sprintf("%s:%s", p.appDetail.Repo, p.appDetail.Version))
	if err != nil {
		log.Error("获取镜像端口失败:", err)
		return err
	}

	if p.appDetail.NginxConfig != "" {
		log.Info("添加Nginx location配置")
		err = p.nm.AddLocation(nginx.NewLocationConfig(p.app.Key, p.containerName).WithTemplate(p.appDetail.NginxConfig).WithPort(port))

		if err != nil {
			log.Error("添加Nginx配置失败:", err)

			std, err := compose.Operate(pluginHelper.GetComposeFile(p.appKey), "stop")
			if err != nil {
				log.Error("停止容器失败:", std, err)
			}
			_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(p.appInstalled.ID)).Update(repo.AppInstalled.Status, model.PluginStatusUpErr)
			return err
		}

		// 提取location
		locations, _ := p.nm.ExtractLocationsByKey(p.app.Key)

		if len(locations) > 0 {
			_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(p.appInstalled.ID)).Update(repo.AppInstalled.Location, locations[0])
		}
	}
	log.Info("Nginx配置完成")
	return nil
}
