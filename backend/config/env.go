package config

import (
	"doo-store/backend/logging"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var EnvConfig = envConfigSchema{}

func (s *envConfigSchema) GetDSN() string {
	return dsn
}

func (s *envConfigSchema) GetGormDialector() gorm.Dialector {
	switch s.STORAGE {
	case "mysql":
		return mysql.Open(s.GetDSN())
	default:
		return sqlite.Open(s.SQLITE_PATH)
	}

}

func (s *envConfigSchema) GetNginxContainerName() string {
	return fmt.Sprintf("dootask-nginx-%s", s.APP_ID)
}

func (s *envConfigSchema) GetNetworkName() string {
	return fmt.Sprintf("dootask-networks-%s", s.APP_ID)
}

func (s *envConfigSchema) GetDefaultContainerName(key string) string {
	return fmt.Sprintf("dootask-plugin-%s-%s-%s", key, key, s.APP_ID)
}

var dsn string

func init() {
	envInit()
	systemInit()
	dsn = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", EnvConfig.MYSQL_USERNAME, EnvConfig.MYSQL_PASSWORD, EnvConfig.MYSQL_HOST, EnvConfig.MYSQL_PORT, EnvConfig.MYSQL_DB_NAME)
}

var defaultConfig = envConfigSchema{
	ENV: "dev",

	APP_ID: "",

	STORAGE: "sqlite",

	SQLITE_PATH: "./app.db",

	MYSQL_HOST:     "127.0.0.1",
	MYSQL_PORT:     "18888",
	MYSQL_USERNAME: "devlop",
	MYSQL_PASSWORD: "123456",
	MYSQL_DB_NAME:  "devlop",

	DATA_DIR: "",

	APP_PREFIX: "dootask-plugin-",

	EXTERNAL_NETWORK_NAME:    "",
	EXTERNAL_NETWORK_IP:      "",
	EXTERNAL_NETWORK_GATEWAY: "",

	YoudaoAppKey:    "",
	YoudaoAppSecret: "",

	DB_PREFIX: "",
}

type envConfigSchema struct {
	ENV string `env:"ENV,DREAM_ENV"`

	APP_ID string

	STORAGE string

	SQLITE_PATH string

	MYSQL_HOST     string
	MYSQL_PORT     string
	MYSQL_USERNAME string
	MYSQL_PASSWORD string
	MYSQL_DB_NAME  string

	DATA_DIR string

	APP_PREFIX string

	EXTERNAL_NETWORK_NAME    string
	EXTERNAL_NETWORK_IP      string
	EXTERNAL_NETWORK_GATEWAY string

	YoudaoAppKey    string
	YoudaoAppSecret string

	DB_PREFIX string
}

func (s *envConfigSchema) IsDev() bool {
	return s.ENV == "dev" || s.ENV == "TESTING"
}

// envInit Reads .env as environment variables and fill corresponding fields into EnvConfig.
// To use a value from EnvConfig , simply call EnvConfig.FIELD like EnvConfig.ENV
// Note: Please keep Env as the first field of envConfigSchema for better logging.
func envInit() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file, ignored")
	}
	v := reflect.ValueOf(defaultConfig)
	typeOfV := v.Type()

	for i := 0; i < v.NumField(); i++ {
		envNameAlt := make([]string, 0)
		fieldName := typeOfV.Field(i).Name
		fieldType := typeOfV.Field(i).Type
		fieldValue := v.Field(i).Interface()

		envNameAlt = append(envNameAlt, fieldName)
		if fieldTag, ok := typeOfV.Field(i).Tag.Lookup("env"); ok && len(fieldTag) > 0 {
			tags := strings.Split(fieldTag, ",")
			envNameAlt = append(envNameAlt, tags...)
		}

		switch fieldType {
		case reflect.TypeOf(0):
			{
				configDefaultValue, ok := fieldValue.(int)
				if !ok {
					logging.Logger.WithFields(logrus.Fields{
						"field": fieldName,
						"type":  "int",
						"value": fieldValue,
						"env":   envNameAlt,
					}).Warningf("Failed to parse default value")
					continue
				}
				envValue := resolveEnv(envNameAlt, fmt.Sprintf("%d", configDefaultValue))
				if EnvConfig.IsDev() {
					fmt.Printf("Reading field[ %s ] default: %v env: %s\n", fieldName, configDefaultValue, envValue)
				}
				if len(envValue) > 0 {
					envValueInteger, err := strconv.ParseInt(envValue, 10, 64)
					if err != nil {
						logging.Logger.WithFields(logrus.Fields{
							"field": fieldName,
							"type":  "int",
							"value": fieldValue,
							"env":   envNameAlt,
						}).Warningf("Failed to parse env value, ignored")
						continue
					}
					reflect.ValueOf(&EnvConfig).Elem().Field(i).SetInt(envValueInteger)
				}
				continue
			}
		case reflect.TypeOf(""):
			{
				configDefaultValue, ok := fieldValue.(string)
				if !ok {
					logging.Logger.WithFields(logrus.Fields{
						"field": fieldName,
						"type":  "int",
						"value": fieldValue,
						"env":   envNameAlt,
					}).Warningf("Failed to parse default value")
					continue
				}
				envValue := resolveEnv(envNameAlt, configDefaultValue)

				if EnvConfig.IsDev() {
					fmt.Printf("Reading field[ %s ] default: %v env: %s\n", fieldName, configDefaultValue, envValue)
				}
				if len(envValue) > 0 {
					reflect.ValueOf(&EnvConfig).Elem().Field(i).SetString(envValue)
				}
			}
		}

	}
}

func resolveEnv(configKeys []string, defaultValue string) string {
	for _, item := range configKeys {
		envValue := os.Getenv(item)
		if envValue != "" {
			return envValue
		}
	}
	return defaultValue
}
