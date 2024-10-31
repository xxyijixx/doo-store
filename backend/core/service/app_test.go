package service

import (
	"doo-store/backend/constant"
	"doo-store/backend/core/model"
	"doo-store/backend/core/repo"
	"testing"
)

func TestApp(t *testing.T) {
	app := &model.App{
		Name:   "Redis",
		Key:    "redis",
		Type:   "databases",
		Status: constant.AppNormal,
	}
	err := repo.App.Create(app)
	if err != nil {
		t.Error(err)
	}
	repo.AppDetail.Create(&model.AppDetail{
		AppID:   app.ID,
		Version: "7.41",
		Params:  "{\"formFields\":[{\"default\":\"redis\",\"envKey\":\"PANEL_REDIS_ROOT_PASSWORD\",\"labelEn\":\"Password\",\"labelZh\":\"密码\",\"random\":true,\"required\":false,\"rule\":\"paramComplexity\",\"type\":\"password\"},{\"default\":6379,\"envKey\":\"PANEL_APP_PORT_HTTP\",\"labelEn\":\"Port\",\"labelZh\":\"端口\",\"required\":true,\"rule\":\"paramPort\",\"type\":\"number\"}]}",
		DockerCompose: `services:
  redis:
    image: redis:7.4.1
    restart: always
    container_name: ${CONTAINER_NAME}
    networks:
      - doo-store-app-network
    ports:
      - 6379:6379
    volumes:
      - ./data:/data
      - ./logs:/logs
    labels:
      createdBy: "Apps"
networks:
  doo-store-app-network:
    external: true`,
		Status: constant.AppNormal,
	})
}
