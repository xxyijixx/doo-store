package service

import (
	"bufio"
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
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"gorm.io/gorm"
)

type AppService struct {
}

type IAppService interface {
	AppPage(ctx dto.ServiceContext, req request.AppSearch) (*dto.PageResult, error)
	AppDetailByKey(ctx dto.ServiceContext, key string) (*response.AppDetail, error)
	AppInstall(ctx dto.ServiceContext, req request.AppInstall) error
	AppInstallOperate(ctx dto.ServiceContext, req request.AppInstalledOperate) error
	AppUnInstall(ctx dto.ServiceContext, req request.AppUnInstall) error
	AppInstalledPage(ctx dto.ServiceContext, req request.AppInstalledSearch) (*dto.PageResult, error)
	Params(ctx dto.ServiceContext, id int64) (any, error)
	UpdateParams(ctx dto.ServiceContext, req request.AppInstall) (any, error)
	AppTags(ctx dto.ServiceContext) ([]*model.Tag, error)
	GetLogs(ctx dto.ServiceContext, req request.AppLogsSearch) (any, error)
}

func NewIAppService() IAppService {
	return &AppService{}
}

func (*AppService) AppPage(ctx dto.ServiceContext, req request.AppSearch) (*dto.PageResult, error) {
	var query repo.IAppDo
	query = repo.App.Order(repo.App.Sort.Desc())
	if req.Name != "" {
		query = query.Where(repo.App.Name.Like(fmt.Sprintf("%%%s%%", req.Name)))
	}
	if req.Class != "" {
		query = query.Where(repo.App.Class.Eq(req.Class))
	}
	if req.ID != 0 {
		query = query.Where(repo.App.ID.Eq(req.ID))
	}
	if req.Description != "" {
		query = query.Where(repo.App.Description.Like(fmt.Sprintf("%%%s%%", req.Description)))
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

func (*AppService) AppDetailByKey(ctx dto.ServiceContext, key string) (*response.AppDetail, error) {

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

func (*AppService) AppInstall(ctx dto.ServiceContext, req request.AppInstall) error {
	app, err := repo.App.Where(repo.App.Key.Eq(req.Key)).First()
	if err != nil {
		log.Info("Error query app")
		return err
	}

	// 检测版本
	dootaskService := NewIDootaskService()
	versionInfoResp, err := dootaskService.GetVersoinInfo()
	if err != nil {
		return err
	}
	check, err := versionInfoResp.CheckVersion(app.DependsVersion)
	if err != nil {
		return err
	}
	if !check {
		// return fmt.Errorf("当前版本不满足要求，需要版本%s以上", app.DependsVersion)
		return e.WithMap(ctx.C, constant.ErrPluginVersionNotSupport, map[string]interface{}{
			"detail": app.DependsVersion,
		}, nil)
	}

	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.AppID.Eq(app.ID)).First()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("安装失败")
		}
	}
	if appInstalled != nil {
		return errors.New("无需重复安装")
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.AppID.Eq(app.ID)).First()
	if err != nil {
		log.Info("Error query app detail")
		return err
	}

	// 检测 docker-compose 文件
	err = compose.Check(req.DockerCompose)
	if err != nil {
		log.Info("DockerCompose 内容未通过检测", err)
		return err
	}

	appKey := config.EnvConfig.APP_PREFIX + app.Key
	// 创建工作目录
	workspaceDir := path.Join(constant.AppInstallDir, appKey)
	err = createDir(workspaceDir)
	if err != nil {
		log.Info("Error create dir", err)
		return err
	}
	// 名称
	name := fmt.Sprintf("plugin-%s", common.RandString(6))
	containerName := config.EnvConfig.APP_PREFIX + app.Key + "-" + name

	paramJson, err := json.Marshal(req.Params)
	if err != nil {
		return err
	}

	// 资源限制
	req.Params[constant.CPUS] = req.CPUS
	req.Params[constant.MemoryLimit] = req.MemoryLimit

	envContent, envJson, err := docker.GenEnv(appKey, containerName, req.Params, false)
	if err != nil {
		return err
	}
	appInstalled = &model.AppInstalled{
		Name:          name,
		AppID:         app.ID,
		AppDetailID:   appDetail.ID,
		Class:         app.Class,
		Repo:          appDetail.Repo,
		Version:       appDetail.Version,
		Params:        string(paramJson),
		Env:           envJson,
		DockerCompose: req.DockerCompose,
		Key:           app.Key,
		Status:        constant.Installing,
	}
	err = appUp(appInstalled, envContent)
	if err != nil {
		log.Info("启动失败", err)
		return err
	}

	// 添加Nginx配置
	client, err := docker.NewClient()
	if err != nil {
		return err
	}
	port, err := client.GetImageFirstExposedPortByName(fmt.Sprintf("%s:%s", appDetail.Repo, appDetail.Version))
	if err != nil {
		return err
	}
	if port != 0 {
		nginx.AddLocation(app.Key, containerName, port)
	}

	return nil
}

func (*AppService) AppInstallOperate(ctx dto.ServiceContext, req request.AppInstalledOperate) error {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.Key.Eq(req.Key)).First()
	if err != nil {
		return err
	}
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := docker.GetComposeFile(appKey)

	supportActions := []string{"start", "stop"}
	if !common.InArray(req.Action, supportActions) {
		return errors.New("不支持的action")
	}

	if req.Action == "stop" {
		err := appStop(appInstalled)
		return err
	}

	stdout, err := compose.Operate(composeFile, req.Action)
	if err != nil {
		log.Info("Error docker compose operate")
		return err
	}
	fmt.Println(stdout)
	insertLog(appInstalled.ID, stdout)
	return nil
}

func (*AppService) AppUnInstall(ctx dto.ServiceContext, req request.AppUnInstall) error {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.Key.Eq(req.Key)).First()
	if err != nil {
		return err
	}
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := docker.GetComposeFile(appKey)
	err = repo.DB.Transaction(func(tx *gorm.DB) error {
		_, err = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Delete()
		if err != nil {
			log.Info("删除插件失败", err)
			return err
		}
		_, err = repo.Use(tx).App.Where(repo.App.ID.Eq(appInstalled.AppID)).Update(repo.App.Status, constant.AppUnused)
		if err != nil {
			log.Info("更新插件状态失败", err)
			return err
		}
		stdout, err := compose.Down(composeFile)
		if err != nil {
			log.Info("Error docker compose down")
			return err
		}
		fmt.Println(stdout)
		return err
	})
	if err != nil {
		log.Info("插件卸载失败", err)
		return errors.New("插件卸载失败")
	}

	nginx.RemoveLocation(appInstalled.Key)
	// 删除compose目录
	_ = os.RemoveAll(fmt.Sprintf("%s/%s", constant.AppInstallDir, appKey))

	return nil
}

func (*AppService) AppInstalledPage(ctx dto.ServiceContext, req request.AppInstalledSearch) (*dto.PageResult, error) {

	query := repo.AppInstalled.Join(repo.App, repo.App.ID.EqCol(repo.AppInstalled.AppID))
	if req.Class != "" {
		query = query.Where(repo.AppInstalled.Class.Eq(req.Class))
	}
	if req.Name != "" {
		query = query.Where(repo.App.Name.Like(fmt.Sprintf("%%%s%%", req.Name)))
	}
	if req.Description != "" {
		query = query.Where(repo.App.Description.Like(fmt.Sprintf("%%%s%%", req.Description)))
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
		return nil, errors.New("查询已安装插件失败")
	}

	pageResult := &dto.PageResult{
		Total: count,
		Items: result,
	}
	return pageResult, nil
}

func (*AppService) Params(ctx dto.ServiceContext, id int64) (any, error) {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(id)).First()
	if err != nil {
		log.Info("Error query app installed", err)
		return nil, errors.New("获取安装插件信息失败")
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.ID.Eq(appInstalled.AppDetailID)).First()
	if err != nil {
		log.Info("Error query app detail", err)
		return nil, errors.New("获取安装插件信息失败")
	}
	// appDetail.Params
	// 解析原始参数
	params := response.AppParams{}
	err = common.StrToStruct(appDetail.Params, &params)
	if err != nil {
		log.Info("错误解析Json", err)
		return nil, err
	}
	env := map[string]string{}
	err = json.Unmarshal([]byte(appInstalled.Env), &env)
	if err != nil {
		log.Info("解析环境变量失败", err)
		return nil, err
	}
	for _, formField := range params.FormFields {
		formField.Value = env[formField.EnvKey]
		formField.Key = formField.EnvKey
	}
	// 构建插件参数
	aParams := response.AppInstalledParamsResp{
		Params:        params.FormFields,
		DockerCompose: appInstalled.DockerCompose,
		CPUS:          env[constant.CPUS],
		MemoryLimit:   env[constant.MemoryLimit],
	}
	return aParams, nil
}

func (*AppService) UpdateParams(ctx dto.ServiceContext, req request.AppInstall) (any, error) {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(req.InstalledId)).First()
	if err != nil {
		log.Info("Error query app installed", err)
		return nil, errors.New("获取安装插件信息失败")
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.ID.Eq(appInstalled.AppDetailID)).First()
	if err != nil {
		log.Info("Error query app detail", err)
		return nil, errors.New("获取安装插件信息失败")
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
	containerName := config.EnvConfig.APP_PREFIX + appInstalled.Key + "-" + appInstalled.Name

	req.Params[constant.CPUS] = req.CPUS
	req.Params[constant.MemoryLimit] = req.MemoryLimit

	envContent, envJson, err := docker.GenEnv(appKey, containerName, req.Params, false)
	if err != nil {
		log.Info("错误生成环境变量文件", err)
		return nil, errors.New("修改参数失败")
	}
	appInstalled.Env = envJson
	paramJson, err := json.Marshal(req.Params)
	if err != nil {
		return nil, errors.New("解析参数失败")
	}
	appInstalled.Params = string(paramJson)
	_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(appInstalled)
	err = appRe(appInstalled, envContent)
	if err != nil {
		log.Info("重启失败", err)
		return nil, errors.New("插件重启失败")
	}
	// 返回修改后的参数
	env := map[string]string{}
	err = json.Unmarshal([]byte(appInstalled.Env), &env)
	if err != nil {
		log.Info("解析环境变量失败", err)
		return nil, err
	}
	for _, formField := range params.FormFields {
		formField.Value = env[formField.EnvKey]
		formField.Key = formField.EnvKey
	}
	aParams := response.AppInstalledParamsResp{
		Params:        params.FormFields,
		DockerCompose: appInstalled.DockerCompose,
		CPUS:          req.CPUS,
		MemoryLimit:   req.MemoryLimit,
	}
	return aParams, nil
}

func (*AppService) AppTags(ctx dto.ServiceContext) ([]*model.Tag, error) {
	tags, err := repo.Tag.Find()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*model.Tag{}, nil
		}
		return nil, err
	}
	return tags, nil
}

func (*AppService) GetLogs(ctx dto.ServiceContext, req request.AppLogsSearch) (any, error) {
	log.Info("获取日志")
	client, err := docker.NewDockerClient()
	if err != nil {
		log.Info("获取Docker失败", err)
		return nil, err
	}
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(req.Id)).First()
	if err != nil {
		log.Info("Error query app installed", err)
		return nil, errors.New("获取安装插件信息失败")
	}

	if appInstalled.Status != constant.Running {
		return nil, errors.New("插件未运行")
	}

	containerName := config.EnvConfig.APP_PREFIX + appInstalled.Key + "-" + appInstalled.Name
	reader, err := client.ContainerLogs(context.Background(), containerName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		// Follow:     true,
		Since: req.Since,
		Until: req.Until,
		Tail:  fmt.Sprintf("%d", req.Tail),
	})
	if err != nil {
		log.Info("Error query container log", err)
		return nil, errors.New("获取日志失败")
	}
	sb := strings.Builder{}

	scanner := bufio.NewScanner(reader)

	cleanText := func(input string) string {
		// 去除所有不可见的控制字符 (ASCII 0-31 和 127)
		reControlChars := regexp.MustCompile(`[\x00-\x1F\x7F]`)
		return reControlChars.ReplaceAllString(input, "")
	}
	date := time.Now()
	go func() {

		for scanner.Scan() {
			fmt.Println(scanner.Text())
			_, err = sb.WriteString(cleanText(scanner.Text()) + "\n")
			if err != nil {
				log.Info("写入字符串时失败", err)
			}
			date = time.Now()
		}

		if err := scanner.Err(); err != nil {
			log.Info("Error scan container log", err)
			return
		}

	}()
	for {
		time.Sleep(time.Microsecond * 500)
		if time.Since(date) > time.Second {
			log.Info("距离上次读取时间超过1秒，读取结束")
			break
		}
	}
	return sb.String(), nil
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
		log.Info("Error WriteFile", err)
		return err
	}
	_, err = docker.WrietEnvFile(appKey, envContent)
	if err != nil {
		log.Info("Error WriteFile", err)
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
		_, err := repo.Use(tx).App.Where(repo.App.ID.Eq(appInstalled.AppID)).Update(repo.App.Status, constant.AppInUse)
		if err != nil {
			return err
		}
		err = repo.Use(tx).AppInstalled.Create(appInstalled)
		if err != nil {
			return err
		}
		composeFile, err := docker.WriteComposeFile(appKey, appInstalled.DockerCompose)
		log.Info("Docker容器UP,", composeFile)
		if err != nil {
			log.Info("Error WriteFile", err)
			return err
		}
		_, err = docker.WrietEnvFile(appKey, envContent)
		if err != nil {
			log.Info("Error WriteFile", err)
			return err
		}
		stdout, err := compose.Up(composeFile)
		if err != nil {
			stdout, err = docker.ParseError(stdout, err)
			log.Info("Error docker compose up:", stdout, err)
			_, _ = repo.Use(tx).AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.UpErr)
			return err
		}
		_, _ = repo.Use(tx).AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Running)
		fmt.Println(stdout)

		insertLog(appInstalled.ID, stdout)
		return nil
	})
	if err != nil {
		insertLog(appInstalled.ID, err.Error())
	} else {
		insertLog(appInstalled.ID, "插件启动")
	}
	return err
}

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
	insertLog(appInstalled.ID, stdout)
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

func insertLog(appInstalledId int64, content string) {
	if content == "" {
		log.Info("log content is empty")
		return
	}
	err := repo.AppLog.Create(&model.AppLog{
		AppInstalledId: appInstalledId,
		Content:        content,
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
