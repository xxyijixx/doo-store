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
	"doo-store/backend/task"
	"doo-store/backend/utils/common"
	"doo-store/backend/utils/compose"
	"doo-store/backend/utils/docker"
	"doo-store/backend/utils/nginx"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"gorm.io/gorm"
)

type AppService struct {
}

type IAppService interface {
	ListApps(ctx dto.ServiceContext, req request.AppSearch) (*dto.PageResult, error)
	GetAppDetail(ctx dto.ServiceContext, key string) (*response.AppDetail, error)
	InstallApp(ctx dto.ServiceContext, req request.AppInstall) error
	UpdateAppInstall(ctx dto.ServiceContext, req request.AppInstalledOperate) error
	UninstallApp(ctx dto.ServiceContext, req request.AppUnInstall) error
	ListInstalledApps(ctx dto.ServiceContext, req request.AppInstalledSearch) (*dto.PageResult, error)
	GetAppParams(ctx dto.ServiceContext, id int64) (any, error)
	UpdateAppParams(ctx dto.ServiceContext, req request.AppInstall) (any, error)
	ListAppTags(ctx dto.ServiceContext) ([]*model.Tag, error)
	GetAppLogs(ctx dto.ServiceContext, req request.AppLogsSearch) (any, error)
	UploadApp(ctx dto.ServiceContext, req request.PluginUpload) error
	GetInstalledAppInfo(ctx dto.ServiceContext, req request.GetInstalledPluginInfo) (*response.GetInstalledPluginInfoResp, error)
	ListRunningAppKeys(ctx dto.ServiceContext) (any, error)
}

func NewIAppService() IAppService {
	return &AppService{}
}

func (*AppService) ListApps(ctx dto.ServiceContext, req request.AppSearch) (*dto.PageResult, error) {
	var query repo.IAppDo
	query = repo.App.Order(repo.App.Sort.Desc())
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 9
	} else if req.PageSize > 1000 {
		req.PageSize = 1000
	}
	if req.Class != "" {
		query = query.Where(repo.App.Class.Eq(req.Class))
	}
	if req.ID != 0 {
		query = query.Where(repo.App.ID.Eq(req.ID))
	}
	if req.Name != "" || req.Description != "" {
		query = query.Where(repo.App.Name.Like(fmt.Sprintf("%%%s%%", req.Name))).Or(repo.App.Description.Like(fmt.Sprintf("%%%s%%", req.Description)))
	}
	result, count, err := query.FindByPage((req.Page-1)*req.PageSize, req.PageSize)

	if err != nil {
		return nil, err
	}

	pageResult := &dto.PageResult{
		Total: count,
		Items: result,
	}
	return pageResult, nil
}

func (*AppService) GetAppDetail(ctx dto.ServiceContext, key string) (*response.AppDetail, error) {

	app, err := repo.App.Where(repo.App.Key.Eq(key)).First()
	if err != nil {
		return nil, err
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.AppID.Eq(app.ID)).First()
	if err != nil {
		return nil, err
	}
	params := response.AppParams{}
	err = common.StrToStruct(appDetail.Params, &params)
	if err != nil {
		return nil, err
	}
	resp := &response.AppDetail{
		AppDetail: *appDetail,
		Params:    params,
	}

	return resp, nil
}

// InstallApp 插件安装
func (*AppService) InstallApp(ctx dto.ServiceContext, req request.AppInstall) error {
	appInstallProcess := NewAppInstallProcess(ctx, req)

	if err := appInstallProcess.ValidateInstallRequirements(); err != nil {
		return err
	}
	if err := appInstallProcess.DHCP(); err != nil {
		return err
	}
	if err := appInstallProcess.ValidateParam(); err != nil {
		return err
	}
	// 异步处理
	manager := task.GetAsyncTaskManager()
	manager.AddTask(func() error {
		if err := appInstallProcess.Install(); err != nil {
			return err
		}
		if err := appInstallProcess.AddNginx(); err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (*AppService) UpdateAppInstall(ctx dto.ServiceContext, req request.AppInstalledOperate) error {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.Key.Eq(req.Key)).First()
	if err != nil {
		return err
	}
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := docker.GetComposeFile(appKey)

	supportActions := []string{"start", "stop"}
	if !common.InArray(req.Action, supportActions) {
		return errors.New(constant.ErrPluginUnsupportedAction)
	}

	if req.Action == "stop" {
		err := appStop(appInstalled)
		return err
	}
	stdout := ""
	if req.Action == "start" {
		// 插件未正常启动，执行up操作
		if appInstalled.Status == constant.UpErr {
			stdout, err = compose.Up(composeFile)
		} else {
			stdout, err = compose.Operate(composeFile, req.Action)
		}
		if err != nil {
			log.Info("Error docker compose operate")

			_, err = docker.ParseError(stdout, err)
			_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
				map[string]interface{}{
					repo.AppInstalled.Message.ColumnName().String(): err.Error(),
				},
			)
			return err
		}
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
			map[string]interface{}{
				repo.AppInstalled.Status.ColumnName().String():  constant.Running,
				repo.AppInstalled.Message.ColumnName().String(): "",
			},
		)
	}

	insertLog(appInstalled.ID, fmt.Sprintf("插件操作[%s]", req.Action), stdout)
	return nil
}

func (*AppService) UninstallApp(ctx dto.ServiceContext, req request.AppUnInstall) error {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.Key.Eq(req.Key)).First()
	if err != nil {
		return err
	}
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := docker.GetComposeFile(appKey)
	err = repo.DB.Transaction(func(tx *gorm.DB) error {
		_, err = repo.Use(tx).AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Delete()
		if err != nil {
			log.Info("删除插件失败", err)
			return err
		}
		_, err = repo.Use(tx).App.Where(repo.App.ID.Eq(appInstalled.AppID)).Update(repo.App.Status, constant.AppUnused)
		if err != nil {
			log.Info("更新插件状态失败", err)
			return err
		}
		if appInstalled.Status != constant.UpErr {
			stdout, err := compose.Down(composeFile)
			if err != nil {
				log.Info("Error docker compose down")
				return err
			}
			fmt.Println(stdout)
		}
		return err
	})
	if err != nil {
		log.Info("插件卸载失败", err)
		return errors.New(constant.ErrPluginUninstallFailed)
	}
	// 释放IP
	docker.GlobalIPAllocator.ReleaseIP(appInstalled.IpAddress)
	nm, err := nginx.NewNginxManager()
	if err != nil {
		return err
	}
	nm.RemoveLocation(appInstalled.Key)
	// 删除compose目录
	_ = os.RemoveAll(fmt.Sprintf("%s/%s", constant.AppInstallDir, appKey))

	return nil
}

func (*AppService) ListInstalledApps(ctx dto.ServiceContext, req request.AppInstalledSearch) (*dto.PageResult, error) {

	query := repo.AppInstalled.Join(repo.App, repo.App.ID.EqCol(repo.AppInstalled.AppID))
	if req.Class != "" {
		query = query.Where(repo.AppInstalled.Class.Eq(req.Class))
	}
	if req.Description != "" || req.Name != "" {
		query = query.Where(repo.App.Description.Like(fmt.Sprintf("%%%s%%", req.Description))).Or(repo.App.Name.Like(fmt.Sprintf("%%%s%%", req.Name)))
	}

	result := []map[string]any{}
	count, err := query.Select(repo.AppInstalled.ALL, repo.App.Icon, repo.App.Description, repo.App.Name).ScanByPage(&result, (req.Page-1)*req.PageSize, req.PageSize)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &dto.PageResult{
				Total: 0,
				Items: []interface{}{},
			}, nil
		}
		log.Info("查询已安装插件失败", err)
		return nil, errors.New(constant.ErrPluginInfoFailed)
	}

	pageResult := &dto.PageResult{
		Total: count,
		Items: result,
	}
	return pageResult, nil
}

func (*AppService) GetAppParams(ctx dto.ServiceContext, id int64) (any, error) {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(id)).First()
	if err != nil {
		log.Info("Error query app installed", err)
		return nil, errors.New(constant.ErrPluginInfoFailed)
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.ID.Eq(appInstalled.AppDetailID)).First()
	if err != nil {
		log.Info("Error query app detail", err)
		return nil, errors.New(constant.ErrPluginInfoFailed)
	}
	// appDetail.Params
	// 解析原始参数
	params := response.AppParams{}
	err = common.StrToStruct(appDetail.Params, &params)
	if err != nil {
		log.Info("错误解析Json", err)
		return nil, err
	}
	env := map[string]interface{}{}
	err = json.Unmarshal([]byte(appInstalled.Env), &env)
	if err != nil {
		log.Info("解析环境变量失败", err)
		return nil, err
	}
	// for _, formField := range params.FormFields {
	// 	formField.Value = env[formField.EnvKey]
	// 	formField.Key = formField.EnvKey
	// }
	params.FormFields = dto.FillAndValidateForm(params.FormFields, env)
	// 构建插件参数
	aParams := response.AppInstalledParamsResp{
		Params:        params.FormFields,
		DockerCompose: appInstalled.DockerCompose,
		CPUS:          env[constant.CPUS].(string),
		MemoryLimit:   env[constant.MemoryLimit].(string),
	}
	return aParams, nil
}

func (*AppService) UpdateAppParams(ctx dto.ServiceContext, req request.AppInstall) (any, error) {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(req.InstalledId)).First()
	if err != nil {
		log.Info("Error query app installed", err)
		return nil, errors.New(constant.ErrPluginInfoFailed)
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.ID.Eq(appInstalled.AppDetailID)).First()
	if err != nil {
		log.Info("Error query app detail", err)
		return nil, errors.New(constant.ErrPluginInfoFailed)
	}
	// appDetail.Params
	// 解析原始参数
	params := response.AppParams{}
	err = common.StrToStruct(appDetail.Params, &params)
	if err != nil {
		log.Info("错误解析Json", err)
		return nil, err
	}
	// TODO 参数校验
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	containerName := appInstalled.Name
	ipAddress := appInstalled.IpAddress

	req.Params[constant.CPUS] = req.CPUS

	if req.MemoryLimit == "" || req.MemoryLimit == "0" {
		req.Params[constant.MemoryLimit] = "0"
	} else {
		if req.MemoryUnit != "m" && req.MemoryUnit != "0" {
			req.MemoryUnit = "m"
		}
		req.Params[constant.MemoryLimit] = req.MemoryLimit + req.MemoryUnit
	}

	// TODO 参数更新
	envContent, envJson, err := docker.GenEnv(appKey, containerName, ipAddress, req.Params, false)
	if err != nil {
		log.Info("错误生成环境变量文件", err)
		return nil, errors.New(constant.ErrPluginModifyParamFailed)
	}

	appInstalled.Env = envJson
	paramJson, err := json.Marshal(req.Params)
	if err != nil {
		return nil, errors.New(constant.ErrPluginParamParseFailed)
	}
	appInstalled.Params = string(paramJson)
	_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(appInstalled)
	err = appRe(appInstalled, envContent)
	if err != nil {
		log.Info("重启失败", err)
		insertLog(appInstalled.ID, "插件重启", err.Error())
		return nil, errors.New(constant.ErrPluginRestartFailed)
	}
	// 返回修改后的参数
	env := map[string]interface{}{}
	err = json.Unmarshal([]byte(appInstalled.Env), &env)
	if err != nil {
		log.Info("解析环境变量失败", err)
		return nil, err
	}
	// for _, formField := range params.FormFields {
	// 	formField.Value = env[formField.EnvKey]
	// 	formField.Key = formField.EnvKey
	// }

	params.FormFields = dto.FillAndValidateForm(params.FormFields, env)
	aParams := response.AppInstalledParamsResp{
		Params:        params.FormFields,
		DockerCompose: appInstalled.DockerCompose,
		CPUS:          req.CPUS,
		MemoryLimit:   req.MemoryLimit,
	}
	insertLog(appInstalled.ID, "插件参数修改", "")
	return aParams, nil
}

func (*AppService) ListAppTags(ctx dto.ServiceContext) ([]*model.Tag, error) {
	tags, err := repo.Tag.Find()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*model.Tag{}, nil
		}
		return nil, err
	}
	return tags, nil
}

func (*AppService) GetAppLogs(ctx dto.ServiceContext, req request.AppLogsSearch) (any, error) {
	log.Info("获取日志")

	// 获取 Docker 客户端
	client, err := docker.NewDockerClient()
	if err != nil {
		log.Error("获取 Docker 客户端失败", err)
		return nil, err
	}
	defer client.Close()

	// 查询已安装的插件信息
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(req.Id)).First()
	if err != nil {
		log.Error("查询插件安装信息失败", err)
		return nil, errors.New(constant.ErrPluginInfoFailed)
	}

	// 校验插件状态
	if appInstalled.Status != constant.Running {
		return nil, errors.New(constant.ErrPluginNotRunning)
	}

	// 检查容器是否存在
	_, err = client.ContainerInspect(context.Background(), appInstalled.Name)
	if err != nil {
		log.Error("容器不存在", err)
		return nil, errors.New(constant.ErrPluginNotInstalled)
	}

	// 获取容器日志
	reader, err := client.ContainerLogs(context.Background(), appInstalled.Name, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      req.Since,
		Until:      req.Until,
		Tail:       fmt.Sprintf("%d", req.Tail),
		Follow:     false,
	})
	if err != nil {
		log.Error("获取容器日志失败", err)
		return nil, errors.New(constant.ErrLogGetFailed)
	}
	defer reader.Close()

	// 读取所有日志内容
	logBytes, err := io.ReadAll(reader)
	if err != nil {
		log.Error("读取日志内容失败", err)
		return nil, errors.New(constant.ErrLogReadFailed)
	}

	// 将字节转换为字符串
	logContent := string(logBytes)

	// 按行分割日志
	logLines := strings.Split(logContent, "\n")

	// 处理每一行日志
	var builder strings.Builder
	for i, line := range logLines {
		if len(line) > 8 { // docker log格式前8字节为header
			// 跳过header,直接获取日志内容
			if i > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(line[8:])
		}
	}

	result := builder.String()

	return result, nil
}

// UploadApp 插件上传
func (AppService) UploadApp(ctx dto.ServiceContext, req request.PluginUpload) error {
	key := req.Plugin.Key
	count, err := repo.App.Where(repo.App.Key.Eq(key)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(constant.ErrPluginKeyExist)
	}
	err = repo.DB.Transaction(func(tx *gorm.DB) error {

		app := &model.App{
			Name:           req.Plugin.Name,
			Key:            req.Plugin.Key,
			Icon:           req.Plugin.Icon,
			Class:          req.Plugin.Class,
			Description:    req.Plugin.Description,
			DependsVersion: req.Plugin.DependsVersion,
			Status:         constant.AppUnused,
		}
		err := repo.Use(tx).App.Create(app)
		if err != nil {
			log.Debug(err.Error())
			return err
		}
		tag, _ := repo.Tag.Where(repo.Tag.Key.Eq(req.Plugin.Class)).First()
		if tag == nil {
			_ = repo.Use(tx).Tag.Create(&model.Tag{
				Key:  req.Plugin.Class,
				Name: req.Plugin.Class,
			})
		}

		dockerCompose := req.Plugin.GenComposeFile()

		_, err = compose.PreCheck(dockerCompose)
		if err != nil {
			return err
		}

		nginxConfig := req.NginxConfig
		if nginxConfig == "" {
			nginxConfig = req.Plugin.GenNginxConfig()
		}

		appDetail := &model.AppDetail{
			AppID:          app.ID,
			Repo:           req.Plugin.Repo,
			Version:        req.Plugin.Version,
			DependsVersion: req.Plugin.DependsVersion,
			Params:         req.Plugin.GenParams(),
			DockerCompose:  dockerCompose,
			NginxConfig:    nginxConfig,
			Status:         constant.AppNormal,
		}
		err = repo.Use(tx).AppDetail.Create(appDetail)
		if err != nil {
			log.Debug(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (AppService) GetInstalledAppInfo(ctx dto.ServiceContext, req request.GetInstalledPluginInfo) (*response.GetInstalledPluginInfoResp, error) {
	// 获取已安装并正常运行的插件信息
	info, err := repo.AppInstalled.Where(repo.AppInstalled.Key.Eq(req.Key), repo.AppInstalled.Status.Eq(constant.Running)).First()
	if err != nil {
		log.Info("查询插件安装信息失败", err)
		return nil, err
	}
	resp := &response.GetInstalledPluginInfoResp{
		Name:     info.Name,
		Key:      info.Key,
		Status:   info.Status,
		Location: info.Location,
	}

	// 获取云盘的provider
	if req.Key == "doocloudisk" {
		env := map[string]string{}
		err = json.Unmarshal([]byte(info.Env), &env)
		if err != nil {
			log.Info("解析环境变量失败", err)
			return nil, err
		}

		resp.CloudProvider = env["CLOUD_PROVIDER"]
	}
	return resp, err
}

func (AppService) ListRunningAppKeys(ctx dto.ServiceContext) (any, error) {
	result := []string{}
	err := repo.AppInstalled.Select(repo.AppInstalled.Key).Where(repo.AppInstalled.Status.Eq(constant.Running)).Pluck(repo.AppInstalled.Key, &result)
	if err != nil {
		return []string{}, err
	}
	return result, nil
}

func appRe(appInstalled *model.AppInstalled, envContent string) error {
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := docker.GetComposeFile(appKey)
	_, err := compose.Down(composeFile)
	if err != nil {
		log.Info("Error docker compose down", err)
		return err
	}
	_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Installing)
	// 写入docker-compose.yaml和环境文件
	composeFile, err = docker.WriteComposeFile(appKey, appInstalled.DockerCompose)
	if err != nil {
		log.Error("DockerCompose文件写入失败", err)
		return err
	}
	_, err = docker.WriteEnvFile(appKey, envContent)
	if err != nil {
		log.Error("环境变量文件写入失败", err)
		return err
	}
	stdout, err := compose.Up(composeFile)
	if err != nil {
		log.Info("Error docker compose up", stdout)
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.UpErr)
		return err
	}
	_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Running)

	return nil
}

// appUp
// envContent key=value
func appUp(appInstalled *model.AppInstalled, envContent string) error {
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		composeFile, err := docker.WriteComposeFile(appKey, appInstalled.DockerCompose)
		log.Info("Docker容器UP,", composeFile)
		if err != nil {
			log.Info("Error WriteFile", err)
			return err
		}
		_, err = docker.WriteEnvFile(appKey, envContent)
		if err != nil {
			log.Info("Error WriteFile", err)
			return err
		}
		stdout, err := compose.Up(composeFile)
		if err != nil {
			stdout, err = docker.ParseError(stdout, err)
			log.Info("Error docker compose up:", stdout, err)
			return err
		}
		fmt.Println(stdout)
		_, err = repo.Use(tx).AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
			model.AppInstalled{
				Status:  constant.Running,
				Message: "",
			},
		)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
			model.AppInstalled{
				Status:  constant.UpErr,
				Message: err.Error(),
			},
		)
		insertLog(appInstalled.ID, "插件启动", err.Error())
	} else {
		insertLog(appInstalled.ID, "插件启动", "")
	}
	return err
}

// appStop 插件停止
func appStop(appInstalled *model.AppInstalled) error {
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := docker.GetComposeFile(appKey)
	_, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Stopped)
	if err != nil {
		return err
	}
	stdout, err := compose.Stop(composeFile)
	if err != nil {
		return fmt.Errorf("error docker compose stop: %s", err.Error())
	}
	insertLog(appInstalled.ID, "插件停止", stdout)
	return nil
}

func createDir(dirPath string) error {
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		if os.IsExist(err) {
			log.WithField("file", dirPath).Debug("file exists")
			return nil
		}
		return err
	}
	return nil
}

func insertLog(appInstalledId int64, prefix, content string) {
	if prefix == "" && content == "" {
		log.Info("log content is empty")
		return
	}
	err := repo.AppLog.Create(&model.AppLog{
		AppInstalledId: appInstalledId,
		Content:        fmt.Sprintf("%s-%s", prefix, content),
	})
	if err != nil {
		log.Info("Error create app log")
	}
}

// ConvertToUTF8 尝试将非 UTF-8 内容转换为 UTF-8
func ConvertToUTF8(input []byte) (string, error) {
	// 尝试使用 GBK 解码（示例，可以替换为其他编码）
	reader := transform.NewReader(strings.NewReader(string(input)), simplifiedchinese.GBK.NewDecoder())
	converted, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(converted), nil
}
