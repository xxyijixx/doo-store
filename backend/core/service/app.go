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
	"os"
	"path"
	"strings"

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
	fmt.Printf("AppInstallDir: %s, DataDir: %s\n", constant.DataDir, constant.AppInstallDir)

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

	workspaceDir := path.Join(constant.AppInstallDir, app.Key)
	err = createDir(workspaceDir)
	if err != nil {
		log.Debug("Error create dir")
		return err
	}

	containerName := config.EnvConfig.APP_PREFIX + app.Key + "-" + req.Name

	// 生成环境变量文件
	envContent := ""
	envContent += fmt.Sprintf("%s=%s\n", "CONTAINER_NAME", containerName)
	for key, value := range req.Params {
		envContent += fmt.Sprintf("%s=%s\n", key, value)
	}
	paramJson, err := json.Marshal(req.Params)
	if err != nil {
		return err
	}
	envMap := map[string]string{}
	envContentLine := strings.Split(envContent, "\n")
	for _, line := range envContentLine {
		env := strings.Split(line, "=")
		if len(env) != 2 {
			continue
		}
		envMap[env[0]] = env[1]
	}
	jsonEnv, err := json.Marshal(envMap)
	if err != nil {
		return err
	}
	appInstalled := &model.AppInstalled{
		Name:          req.Name,
		AppID:         app.ID,
		AppDetailID:   appDetail.ID,
		Version:       appDetail.Version,
		Params:        string(paramJson),
		Env:           string(jsonEnv),
		DockerCompose: appDetail.DockerCompose,
		Key:           app.Key,
		Status:        constant.Installing,
	}
	repo.AppInstalled.Create(appInstalled)

	composeFile := fmt.Sprintf("%s/%s/docker-compose.yml", constant.AppInstallDir, app.Key)
	err = os.WriteFile(composeFile, []byte(appDetail.DockerCompose), 0644)
	if err != nil {
		log.Debug("Error WriteFile", err)
		return err
	}
	envFile := fmt.Sprintf("%s/%s/.env", constant.AppInstallDir, app.Key)
	err = os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		log.Debug("Error WriteFile", err)
		return err
	}
	stdout, err := compose.Up(composeFile)
	if err != nil {
		log.Debug("Error docker compose up")
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.UpErr)
		return err
	}
	fmt.Println(stdout)
	insertLog(appInstalled.ID, stdout)
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
	composeFile := fmt.Sprintf("%s/%s/docker-compose.yml", constant.AppInstallDir, appInstalled.Key)

	if req.Action == "update" {
		// 重建容器
		_, err := compose.Down(composeFile)
		if err != nil {
			log.Debug("Error docker compose operate")
			return err
		}

		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Stop)

		envFile := fmt.Sprintf("%s/%s/.env", constant.AppInstallDir, appInstalled.Key)
		envContent := ""
		name, exsit := req.Params["name"]
		containerName := config.EnvConfig.APP_PREFIX + appInstalled.Key + "-"
		if exsit {
			containerName += fmt.Sprintf("%s", name)
		} else {
			containerName += appInstalled.Name
		}

		envContent += fmt.Sprintf("%s=%s\n", "CONTAINER_NAME", containerName)
		for key, value := range req.Params {
			envContent += fmt.Sprintf("%s=%s\n", key, value)
		}
		err = os.WriteFile(envFile, []byte(envContent), 0644)
		if err != nil {
			log.Debug("Error WriteFile", err)
			return err
		}
		envMap := map[string]string{}
		envContentLine := strings.Split(envContent, "\n")
		for _, line := range envContentLine {
			env := strings.Split(line, "=")
			if len(env) != 2 {
				continue
			}
			envMap[env[0]] = env[1]
		}
		jsonEnv, err := json.Marshal(envMap)
		if err != nil {
			return err
		}
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Env, string(jsonEnv))
		_, err = compose.Up(composeFile)
		if err != nil {
			_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.UpErr)
			log.Debug("Error docker compose operate")
			return err
		}
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Running)
		return nil
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
	composeFile := fmt.Sprintf("%s/%s/docker-compose.yml", constant.AppInstallDir, appInstalled.Key)
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
	_ = os.RemoveAll(fmt.Sprintf("%s/%s", constant.AppInstallDir, appInstalled.Key))

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
