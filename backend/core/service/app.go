package service

import (
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
	"doo-store/backend/utils/nginx"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AppService struct {
}

type IAppService interface {
	AppPage(req request.AppSearch) (*dto.PageResult, error)
	AppDetailByKey(key string) (*response.AppDetail, error)
	AppInstall(req request.AppInstall) error
	AppInstallOperate(req request.AppInstalledOperate) error
	AppUnInstall(req request.AppUnInstall) error
	AppInstalledPage(req request.AppInstalledSearch) (*dto.PageResult, error)
}

func NewIAppService() IAppService {
	return &AppService{}
}

func (*AppService) AppPage(req request.AppSearch) (*dto.PageResult, error) {

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
	result, count, err := query.FindByPage(req.Page-1, req.PageSize)

	if err != nil {
		return nil, err
	}

	pageResult := &dto.PageResult{
		Total: count,
		Items: result,
	}
	return pageResult, nil
}

func (*AppService) AppDetailByKey(key string) (*response.AppDetail, error) {

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

func (*AppService) AppInstall(req request.AppInstall) error {
	app, err := repo.App.Where(repo.App.Key.Eq(req.Key)).First()
	if err != nil {
		log.Debug("Error query app")
		return err
	}
	_, err = repo.AppInstalled.Where(repo.AppInstalled.AppID.Eq(app.ID)).First()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("无需重复安装")
		}
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.AppID.Eq(app.ID)).First()
	if err != nil {
		log.Debug("Error query app detail")
		return err
	}
	appKey := config.EnvConfig.APP_PREFIX + app.Key
	// 创建工作目录
	workspaceDir := path.Join(constant.AppInstallDir, appKey)
	err = createDir(workspaceDir)
	if err != nil {
		log.Debug("Error create dir")
		return err
	}
	// 如果名称不存在则随机生成
	if req.Name == "" {
		req.Name = fmt.Sprintf("%d", rand.Int31n(100000))
	}
	containerName := config.EnvConfig.APP_PREFIX + app.Key + "-" + req.Name

	paramJson, err := json.Marshal(req.Params)
	if err != nil {
		return err
	}

	envContent, envJson, err := docker.GenEnv(appKey, containerName, req.Params, false)
	if err != nil {
		return err
	}
	appInstalled := &model.AppInstalled{
		Name:          req.Name,
		AppID:         app.ID,
		AppDetailID:   appDetail.ID,
		Version:       appDetail.Version,
		Params:        string(paramJson),
		Env:           envJson,
		DockerCompose: appDetail.DockerCompose,
		Key:           app.Key,
		Status:        constant.Installing,
	}
	err = appUp(appInstalled, envContent)
	if err != nil {
		return err
	}

	// 添加Nginx配置
	client, err := docker.NewClient()
	if err != nil {
		return err
	}
	port, err := client.GetImageFirstExposedPortByName(fmt.Sprintf("%s:%s", app.Key, appDetail.Version))
	if err != nil {
		return err
	}
	if port != 0 {
		nginx.AddLocation(app.Key, containerName, port)
	}

	return nil
}

func (*AppService) AppInstallOperate(req request.AppInstalledOperate) error {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.Key.Eq(req.Key)).First()
	if err != nil {
		return err
	}
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := fmt.Sprintf("%s/%s/docker-compose.yml", constant.AppInstallDir, appKey)

	if req.Action == "update" {
		// 重建容器
		_, err := compose.Down(composeFile)
		if err != nil {
			log.Debug("Error docker compose operate")
			return err
		}

		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Stop)

		name, exsit := req.Params["name"]
		containerName := config.EnvConfig.APP_PREFIX + appInstalled.Key + "-"
		if exsit && name != "" {
			containerName += fmt.Sprintf("%s", name)
		} else {
			containerName += appInstalled.Name
		}
		_, envJson, err := docker.GenEnv(appKey, containerName, req.Params, true)
		if err != nil {
			return err
		}
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Env, envJson)
		_, err = compose.Up(composeFile)
		if err != nil {
			_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.UpErr)
			log.Debug("Error docker compose operate")
			return err
		}
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Running)
		return nil
	}
	if req.Action == "stop" {
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Stopped)
	}
	stdout, err := compose.Operate(composeFile, req.Action)
	if err != nil {
		log.Debug("Error docker compose operate")
		return err
	}
	fmt.Println(stdout)
	insertLog(appInstalled.ID, stdout)
	return nil
}

func (*AppService) AppUnInstall(req request.AppUnInstall) error {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.Key.Eq(req.Key)).First()
	if err != nil {
		return err
	}
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := fmt.Sprintf("%s/%s/docker-compose.yml", constant.AppInstallDir, appKey)
	stdout, err := compose.Down(composeFile)
	if err != nil {
		log.Debug("Error docker compose down")
		return err
	}
	fmt.Println(stdout)
	_, err = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Delete()
	if err != nil {
		return err
	}

	nginx.RemoveLocation(appInstalled.Key)
	// 删除compose目录
	_ = os.RemoveAll(fmt.Sprintf("%s/%s", constant.AppInstallDir, appKey))

	return nil
}

func (*AppService) AppInstalledPage(req request.AppInstalledSearch) (*dto.PageResult, error) {
	var query repo.IAppInstalledDo
	query = repo.AppInstalled.Order(repo.AppInstalled.ID.Desc())

	if req.Class != "" {
		query = query.Where(repo.App.Class.Eq(req.Class))
	}

	result, count, err := query.FindByPage(req.Page-1, req.PageSize)

	if err != nil {
		return nil, err
	}

	pageResult := &dto.PageResult{
		Total: count,
		Items: result,
	}
	return pageResult, nil
}

// appUp
// envContent key=value
func appUp(appInstalled *model.AppInstalled, envContent string) error {
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		_, err := repo.Use(tx).App.Where(repo.App.ID.Eq(appInstalled.ID)).Update(repo.App.Status, constant.AppInUse)
		if err != nil {
			return err
		}
		repo.Use(tx).AppInstalled.Create(appInstalled)

		composeFile := fmt.Sprintf("%s/%s/docker-compose.yml", constant.AppInstallDir, appKey)
		err = os.WriteFile(composeFile, []byte(appInstalled.DockerCompose), 0644)
		if err != nil {
			log.Debug("Error WriteFile", err)
			return err
		}
		envFile := fmt.Sprintf("%s/%s/.env", constant.AppInstallDir, appKey)
		err = os.WriteFile(envFile, []byte(envContent), 0644)
		if err != nil {
			log.Debug("Error WriteFile", err)
			return err
		}
		stdout, err := compose.Up(composeFile)
		if err != nil {
			log.Debug("Error docker compose up")
			insertLog(appInstalled.ID, err.Error())
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
	}
	return err
}

func appStop(appInstalled model.AppInstalled) error {

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
		log.Debug("log content is empty")
		return
	}
	err := repo.AppLog.Create(&model.AppLog{
		AppInstalledId: appInstalledId,
		Content:        content,
	})
	if err != nil {
		log.Debug("Error create app log")
	}
}
