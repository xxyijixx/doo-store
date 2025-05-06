package app

import (
	"doo-store/backend/config"
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
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Debug("File not exist:", filename)
		} else {
			logrus.Debug("Failed to read file", filename)
		}
		return err
	}

	// 解析 JSON
	err = json.Unmarshal(data, &pluginConfig)
	if err != nil {
		logrus.Debug(err.Error())
		return err
	}
	appKeyMap, _ := loadApps()
	oldTagMap, err := loadTags()
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Debug(err.Error())
		return err
	}
	tagMap := make(map[string]struct{})

	err = repo.DB.Transaction(func(tx *gorm.DB) error {
		for _, p := range pluginConfig.Plugins {
			tagMap[p.Class] = struct{}{}
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
				Status:         model.AppUnused,
			}
			err := repo.Use(tx).App.Create(app)
			if err != nil {
				logrus.Debug(err.Error())
				return err
			}
			appKeyMap[p.Key] = struct{}{}

			appDetail := &model.AppDetail{
				AppID:          app.ID,
				Repo:           p.Repo,
				Version:        p.Version,
				DependsVersion: p.DependsVersion,
				Params:         p.GenParams(),
				DockerCompose:  p.GenComposeFile(),
				NginxConfig:    p.GenNginxConfig(),
				Status:         model.AppNormal,
			}
			err = repo.Use(tx).AppDetail.Create(appDetail)
			if err != nil {
				logrus.Debug(err.Error())
				return err
			}
		}
		createTagsIfNeeded(tagMap, oldTagMap, tx)
		return nil
	})
	return err
}

// createTagsIfNeeded 创建缺失的标签
func createTagsIfNeeded(tagMap, oldTagMap map[string]struct{}, tx *gorm.DB) error {
	needTags := make([]*model.Tag, 0)
	for key := range tagMap {
		if _, exist := oldTagMap[key]; !exist {
			needTags = append(needTags, &model.Tag{
				Name: key,
				Key:  key,
			})
		}
	}
	if len(needTags) > 0 {
		if err := repo.Use(tx).Tag.Create(needTags...); err != nil {
			return err
		}
	}
	return nil
}

func loadApps() (map[string]struct{}, error) {
	apps, err := repo.App.Select(repo.App.Key).Find()
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Debug(err, "Failed to find apps")
		return nil, err
	}

	appKeyMap := make(map[string]struct{})
	for _, app := range apps {
		appKeyMap[app.Key] = struct{}{}
	}
	return appKeyMap, nil
}

func loadTags() (map[string]struct{}, error) {
	tags, err := repo.Tag.Find()
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Debug(err.Error())
		return nil, err
	}
	oldTagMap := make(map[string]struct{})
	for _, tag := range tags {
		oldTagMap[tag.Name] = struct{}{}
	}
	return oldTagMap, nil
}
