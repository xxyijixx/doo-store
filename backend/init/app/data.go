package app

import (
	"doo-store/backend/constant"
	"doo-store/backend/core/dto/response"
	"doo-store/backend/core/model"
	"doo-store/backend/core/repo"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Welcome struct {
	Plugins []Plugin `json:"plugins"`
}

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
}

type EnvElement struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
}

type Volume struct {
	Local  string `json:"local"`
	Target string `json:"target"`
}

func LoadData() error {
	logrus.Info("Loading data...")
	var config Welcome
	filename := "./docker/init/data.json"
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		logrus.Debug("File not exist:", filename)
		return err
	}

	if err != nil {
		logrus.Debug(err.Error())
		return err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		logrus.Debug(err.Error())
		return err
	}
	apps, err := repo.App.Find()
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Debug(err.Error())
		return err
	}
	appKeyMap := make(map[string]string)
	for _, app := range apps {
		appKeyMap[app.Key] = "true"
	}
	tags, err := repo.Tag.Find()
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Debug(err.Error())
		return err
	}
	tagMap := make(map[string]string)
	oldTagMap := make(map[string]*model.Tag)
	for _, tag := range tags {
		oldTagMap[tag.Name] = tag
	}
	err = repo.DB.Transaction(func(tx *gorm.DB) error {
		for _, p := range config.Plugins {
			tagMap[p.Class] = "true"
			// 对于key存在，忽略
			if _, exist := appKeyMap[p.Key]; exist {
				continue
			}
			app := &model.App{
				Name:           p.Name,
				Key:            p.Key,
				Icon:           p.Icon,
				Class:          p.Class,
				Description:    p.Description,
				DependsVersion: p.DependsVersion,
				Status:         constant.AppUnused,
			}
			err := repo.Use(tx).App.Create(app)
			if err != nil {
				logrus.Debug(err.Error())
				return err
			}
			appKeyMap[p.Key] = "true"

			appDetail := &model.AppDetail{
				AppID:          app.ID,
				Repo:           p.Repo,
				Version:        p.Version,
				DependsVersion: p.DependsVersion,
				Params:         p.GenParams(),
				DockerCompose:  p.GenComposeFile(),
				Status:         constant.AppNormal,
			}
			err = repo.Use(tx).AppDetail.Create(appDetail)
			if err != nil {
				logrus.Debug(err.Error())
				return err
			}
		}

		needTags := make([]*model.Tag, 0)
		unneedTags := make([]*model.Tag, 0)
		for key := range tagMap {
			if _, exist := oldTagMap[key]; !exist {
				needTags = append(needTags, &model.Tag{
					Name: key,
					Key:  key,
				})
			}
		}
		for key, tag := range oldTagMap {
			if _, exist := tagMap[key]; !exist {
				unneedTags = append(unneedTags, tag)
			}
		}
		repo.Use(tx).Tag.Create(needTags...)
		repo.Use(tx).Tag.Delete(unneedTags...)
		return nil
	})
	//

	return err
}

func (p *Plugin) GenComposeFile() string {
	composeContent := p.GenService()
	composeContent += p.GenNetwork()
	return composeContent
}

func (p *Plugin) GenNetwork() string {
	networkContent := make([]string, 0)
	networkContent = append(networkContent, "networks:")
	networkContent = append(networkContent, fmt.Sprintf("%s%s:", getSpaces(1), constant.AppNetwork))
	networkContent = append(networkContent, fmt.Sprintf("%sexternal: true", getSpaces(2)))

	networkContent = append(networkContent, "\n")

	return strings.Join(networkContent, "\n")
}

func (p *Plugin) GenService() string {
	serviceContent := make([]string, 0)
	serviceContent = append(serviceContent, "services:")
	serviceContent = append(serviceContent, fmt.Sprintf("%s%s:", getSpaces(1), p.Key))
	serviceContent = append(serviceContent, fmt.Sprintf("%simage: %s:%s", getSpaces(2), p.Repo, p.Version))
	serviceContent = append(serviceContent, fmt.Sprintf("%srestart: always", getSpaces(2)))
	serviceContent = append(serviceContent, fmt.Sprintf("%scontainer_name: ${CONTAINER_NAME}", getSpaces(2)))
	// networks:
	serviceContent = append(serviceContent, fmt.Sprintf("%snetworks:", getSpaces(2)))
	serviceContent = append(serviceContent, fmt.Sprintf("%s- %s", getSpaces(3), constant.AppNetwork))

	if len(p.Volume) != 0 {
		serviceContent = append(serviceContent, fmt.Sprintf("%svolumes:", getSpaces(2)))
		for _, v := range p.Volume {
			serviceContent = append(serviceContent, fmt.Sprintf("%s- %s:%s", getSpaces(3), v.Local, v.Target))
		}
	}

	serviceContent = append(serviceContent, fmt.Sprintf("%slabels:", getSpaces(2)))
	serviceContent = append(serviceContent, fmt.Sprintf("%screatedBy: \"Apps\"", getSpaces(3)))

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

func getSpaces(num int) string {
	spaces := strings.Repeat(" ", 2*num)
	return spaces
}
