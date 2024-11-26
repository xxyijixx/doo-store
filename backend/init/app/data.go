package app

import (
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/core/dto"
	"doo-store/backend/core/model"
	"doo-store/backend/core/repo"
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PluginConfig struct {
	Plugins []dto.Plugin `json:"plugins"`
}

func LoadData() error {
	var pluginConfig PluginConfig
	filename := "./docker/init/data.json"
	if config.EnvConfig.ENV == "prod" {
		filename = "./init/data.json"
	}
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		logrus.Debug("File not exist:", filename)
		return err
	}

	if err != nil {
		logrus.Debug(err.Error())
		return err
	}

	err = json.Unmarshal(data, &pluginConfig)
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
		for _, p := range pluginConfig.Plugins {
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
				NginxConfig:    p.GenNginxConfig(),
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
		if len(unneedTags) != 0 {
			repo.Use(tx).Tag.Delete(unneedTags...)
		}
		repo.Use(tx).Tag.Create(needTags...)
		return nil
	})
	return err
}
