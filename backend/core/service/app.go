package service

import (
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/core/dto"
	"doo-store/backend/core/dto/request"
	"doo-store/backend/core/model"
	"doo-store/backend/core/repo"
	"doo-store/backend/utils/compose"
	"doo-store/backend/utils/docker"
	"doo-store/backend/utils/nginx"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

type AppService struct {
}

type IAppService interface {
	AppPage(req request.AppSearch) (*dto.PageResult, error)
	AppDetailByKey(key string) (*model.AppDetail, error)
	AppDetailByKeyAndVersion(key, version string) (*model.AppDetail, error)
	AppInstall(req request.AppInstall) error
}

func NewIAppService() IAppService {
	return &AppService{}
}

func (*AppService) AppPage(req request.AppSearch) (*dto.PageResult, error) {

	result, count, err := repo.App.FindByPage(req.Page, req.PageSize)

	if err != nil {
		return nil, err
	}

	pageResult := &dto.PageResult{
		Total: count,
		Items: result,
	}
	return pageResult, nil
}

func (*AppService) AppDetailByKey(key string) (*model.AppDetail, error) {

	app, err := repo.App.Where(repo.App.Key.Eq(key)).First()
	if err != nil {
		return nil, err
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.AppID.Eq(app.ID)).First()
	if err != nil {
		return nil, err
	}
	return appDetail, nil
}

func (*AppService) AppDetailByKeyAndVersion(key, version string) (*model.AppDetail, error) {
	app, err := repo.App.Where(repo.App.Key.Eq(key)).First()
	if err != nil {
		return nil, err
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.AppID.Eq(app.ID), repo.AppDetail.Version.Eq(version)).First()
	if err != nil {
		return nil, err
	}
	return appDetail, nil
}

func (*AppService) AppInstall(req request.AppInstall) error {
	fmt.Printf("AppInstallDir: %s, DataDir: %s\n", constant.DataDir, constant.AppInstallDir)
	app, err := repo.App.Where(repo.App.Key.Eq(req.Key)).First()
	if err != nil {
		log.Debug("Error query app")
		return err
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.AppID.Eq(app.ID), repo.AppDetail.Version.Eq(req.Version)).First()
	if err != nil {
		log.Debug("Error query app detail")
		return err
	}
	composeContent := appDetail.DockerCompose
	workspaceDir := path.Join(constant.AppInstallDir, app.Key)
	err = createDir(workspaceDir)
	if err != nil {
		log.Debug("Error create dir")
		return err
	}
	workspaceDir = path.Join(workspaceDir, req.Name)
	err = createDir(workspaceDir)
	if err != nil {
		log.Debug("Error create dir")
		return err
	}
	composeFile := fmt.Sprintf("%s/%s/%s/docker-compose.yml", constant.AppInstallDir, app.Key, req.Name)
	containerName := config.EnvConfig.APP_PREFIX + app.Key + "-" + req.Name
	composeContent = strings.ReplaceAll(composeContent, "${CONTAINER_NAME}", containerName)
	for key, value := range req.Params {
		replaceValue := fmt.Sprintf("%v", value)
		composeContent = strings.ReplaceAll(composeContent, fmt.Sprintf("${%s}", key), replaceValue)
	}
	paramJson, err := json.Marshal(req.Params)
	if err != nil {
		return err
	}
	appInstalled := &model.AppInstalled{
		Name:          req.Name,
		AppID:         app.ID,
		AppDetailID:   appDetail.ID,
		Version:       appDetail.Version,
		Params:        string(paramJson),
		DockerCompose: composeContent,
		Status:        constant.Installing,
	}
	repo.AppInstalled.Create(appInstalled)

	err = os.WriteFile(composeFile, []byte(composeContent), 0644)
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
	// 添加Nginx配置
	client, err := docker.NewClient()
	if err != nil {
		return err
	}
	port, err := client.GetExposedFirstPortByName(fmt.Sprintf("%s:%s", app.Key, appDetail.Version))
	if err != nil {
		return err
	}
	if port != 0 {
		nginx.AddLocation(containerName, containerName, port)
	}

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
