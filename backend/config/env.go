package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Language = []string{language.Chinese.String(), language.TraditionalChinese.String(), language.English.String(), language.Korean.String(), language.Japanese.String(), language.German.String(), language.French.String(), language.Indonesian.String()}
)

var EnvConfig = envConfigSchema{}

func (s *envConfigSchema) GetGormDialector() gorm.Dialector {
	switch s.STORAGE {
	case "mysql":
		return mysql.Open(s.MySQL().GetDSN())
	default:
		return sqlite.Open(s.SQLITE_PATH)
	}
}

func (s *envConfigSchema) GetNginxContainerName() string {
	return fmt.Sprintf("dootask-nginx-%s", s.APP_ID)
}

func (s *envConfigSchema) GetDefaultContainerName(key string) string {
	return fmt.Sprintf("dootask-plugin-%s-%s", key, s.APP_ID)
}

func init() {
	envInit()
}

// 应用基本配置
type AppConfig struct {
	ENV                 string
	APP_ID              string
	APP_KEY             string
	STORAGE             string
	SQLITE_PATH         string
	DATA_DIR            string
	PLUGIN_PREFIX       string
	DB_PREFIX           string
	PLUGIN_CIDR         string
	NETWORK_NAME        string
	SHARED_COMPOSE      bool
	SHARED_COMPOSE_NAME string
}

// MySQL数据库配置
type MySQLConfig struct {
	HOST     string
	PORT     string
	USERNAME string
	PASSWORD string
	DB_NAME  string
}

func (s MySQLConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", s.USERNAME, s.PASSWORD, s.HOST, s.PORT, s.DB_NAME)
}

// DooTask相关配置
type DooTaskConfig struct {
	DIR      string
	APP_ID   string
	APP_IPPR string
	APP_KEY  string
	URL      string
}

// DooTask数据库配置
type DooTaskDBConfig struct {
	HOST     string
	PORT     string
	DATABASE string
	USERNAME string
	PASSWORD string
	PREFIX   string
}

// DooTask Redis配置
type DooTaskRedisConfig struct {
	HOST     string
	PORT     string
	PASSWORD string
}

// 第三方服务配置
type ThirdPartyConfig struct {
	YoudaoAppKey    string
	YoudaoAppSecret string
}

// 环境配置结构体
type envConfigSchema struct {
	// 应用基本配置
	ENV                 string
	APP_KEY             string
	APP_ID              string
	APP_IPPR            string
	STORAGE             string
	SQLITE_PATH         string
	DATA_DIR            string
	PLUGIN_PREFIX       string
	DB_PREFIX           string
	PLUGIN_CIDR         string
	SHARED_COMPOSE      bool
	SHARED_COMPOSE_NAME string

	// MySQL数据库配置
	MYSQL_HOST     string
	MYSQL_PORT     string
	MYSQL_USERNAME string
	MYSQL_PASSWORD string
	MYSQL_DB_NAME  string

	// Redis配置
	REDIS_HOST     string
	REDIS_PASSWORD string
	REDIS_PORT     string

	// DooTask相关配置
	DOOTASK_DIR string
	DOOTASK_URL string

	// 第三方服务配置
	YoudaoAppKey    string
	YoudaoAppSecret string
}

// 获取应用配置
func (s *envConfigSchema) App() AppConfig {
	return AppConfig{
		ENV:            s.ENV,
		APP_ID:         s.APP_ID,
		STORAGE:        s.STORAGE,
		SQLITE_PATH:    s.SQLITE_PATH,
		DATA_DIR:       s.DATA_DIR,
		PLUGIN_PREFIX:  s.PLUGIN_PREFIX,
		DB_PREFIX:      s.DB_PREFIX,
		PLUGIN_CIDR:    fmt.Sprintf("%s.30/24", s.APP_IPPR),
		NETWORK_NAME:   fmt.Sprintf("dootask-networks-%s", s.APP_ID),
		SHARED_COMPOSE: s.SHARED_COMPOSE,
		SHARED_COMPOSE_NAME: func() string {
			if s.SHARED_COMPOSE_NAME == "" {
				return s.PLUGIN_PREFIX
			}
			return s.SHARED_COMPOSE_NAME
		}(),
	}
}

// 获取MySQL配置
func (s *envConfigSchema) MySQL() MySQLConfig {
	return MySQLConfig{
		HOST:     s.MYSQL_HOST,
		PORT:     s.MYSQL_PORT,
		USERNAME: s.MYSQL_USERNAME,
		PASSWORD: s.MYSQL_PASSWORD,
		DB_NAME:  s.MYSQL_DB_NAME,
	}
}

// 获取DooTask配置
func (s *envConfigSchema) DooTask() DooTaskConfig {
	dootaskDir := s.DOOTASK_DIR
	if dootaskDir == "" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed to get current directory", err)
		}
		dootaskDir = pwd
	}
	return DooTaskConfig{
		DIR:      dootaskDir,
		APP_ID:   s.APP_ID,
		APP_IPPR: s.APP_IPPR,
		APP_KEY:  s.APP_KEY,
		URL:      s.DOOTASK_URL,
	}
}

// 获取DooTask数据库配置
func (s *envConfigSchema) DooTaskDB() DooTaskDBConfig {
	return DooTaskDBConfig{
		HOST:     s.MYSQL_HOST,
		PORT:     s.MYSQL_PORT,
		DATABASE: s.MYSQL_DB_NAME,
		USERNAME: s.MYSQL_USERNAME,
		PASSWORD: s.MYSQL_PASSWORD,
		PREFIX:   s.DB_PREFIX,
	}
}

// 获取DooTask Redis配置
func (s *envConfigSchema) DooTaskRedis() DooTaskRedisConfig {
	return DooTaskRedisConfig{
		HOST:     s.REDIS_HOST,
		PORT:     s.REDIS_PORT,
		PASSWORD: s.REDIS_PASSWORD,
	}
}

// 获取第三方服务配置
func (s *envConfigSchema) ThirdParty() ThirdPartyConfig {
	return ThirdPartyConfig{
		YoudaoAppKey:    s.YoudaoAppKey,
		YoudaoAppSecret: s.YoudaoAppSecret,
	}
}

func (s *envConfigSchema) IsDev() bool {
	return s.ENV == "dev" || s.ENV == "TESTING"
}

// envInit 使用Viper读取配置并填充到EnvConfig中
// 要使用EnvConfig中的值，只需调用EnvConfig.FIELD，例如EnvConfig.ENV
// 或者使用分组方法，如EnvConfig.App().ENV, EnvConfig.MySQL().HOST等
func envInit() {
	// 初始化Viper
	v := viper.New()

	// 设置默认值
	setDefaults(v)

	// 读取.env文件
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Print("未找到.env文件，将使用默认值和环境变量")
		} else {
			log.Printf("读取.env文件出错: %v，将使用默认值和环境变量", err)
		}
	}

	// 允许从环境变量读取
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置别名，支持多个环境变量名称
	v.RegisterAlias("ENV", "DREAM_ENV")

	// 将配置映射到结构体
	mapConfigToStruct(v)
	// 检查必填项
	requiredFields := []string{"APP_ID", "REDIS_HOST", "REDIS_PORT"}
	for _, field := range requiredFields {
		if v.GetString(field) == "" {
			log.Fatalf("配置项 %s 不能为空", field)
		}
	}
}

// 获取配置值，尝试多个键，如果都没有则返回默认值
func getConfigValue(v *viper.Viper, keys []string, defaultValue string) string {
	for _, key := range keys {
		if v.IsSet(key) {
			value := v.GetString(key)
			if value != "" {
				return value
			}
		}
	}
	return defaultValue
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// 应用基本配置默认值
	v.SetDefault("ENV", "dev")
	v.SetDefault("APP_ID", "")
	v.SetDefault("STORAGE", "sqlite")
	v.SetDefault("SQLITE_PATH", "./app.db")
	v.SetDefault("DATA_DIR", "")
	v.SetDefault("PLUGIN_PREFIX", "plugin-dootask")
	v.SetDefault("PLUGIN_CIDR", "")
	v.SetDefault("SHARED_COMPOSE", true)
	v.SetDefault("SHARED_COMPOSE_NAME", "")

	// MySQL配置默认值
	v.SetDefault("DB_PREFIX", "pre_")

	// Redis配置默认值
	v.SetDefault("REDIS_HOST", "localhost") // Redis配置默认值
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_PORT", "6379")

	// DooTask配置默认值
	v.SetDefault("DOOTASK_DIR", "")
	v.SetDefault("DOOTASK_URL", "http://127.0.0.1:2222")

	// 第三方服务配置默认值
	v.SetDefault("YoudaoAppKey", "")
	v.SetDefault("YoudaoAppSecret", "")
}

// mapConfigToStruct 将Viper配置映射到EnvConfig结构体
func mapConfigToStruct(v *viper.Viper) {
	// 应用基本配置
	EnvConfig.ENV = v.GetString("ENV")
	EnvConfig.APP_ID = v.GetString("APP_ID")
	EnvConfig.APP_IPPR = v.GetString("APP_IPPR")
	EnvConfig.APP_KEY = v.GetString("APP_KEY")
	EnvConfig.STORAGE = v.GetString("STORAGE")
	EnvConfig.SQLITE_PATH = v.GetString("SQLITE_PATH")
	EnvConfig.DATA_DIR = v.GetString("DATA_DIR")
	EnvConfig.PLUGIN_PREFIX = v.GetString("PLUGIN_PREFIX")
	EnvConfig.DB_PREFIX = v.GetString("DB_PREFIX")
	EnvConfig.PLUGIN_CIDR = v.GetString("PLUGIN_CIDR")
	EnvConfig.SHARED_COMPOSE = v.GetBool("SHARED_COMPOSE")
	EnvConfig.SHARED_COMPOSE_NAME = v.GetString("SHARED_COMPOSE_NAME")

	// MySQL配置 - 从多个键获取值
	EnvConfig.MYSQL_HOST = getConfigValue(v, []string{"MYSQL_HOST", "DB_HOST"}, "127.0.0.1")
	EnvConfig.MYSQL_PORT = getConfigValue(v, []string{"MYSQL_PORT", "DB_PORT"}, "18888")
	EnvConfig.MYSQL_USERNAME = getConfigValue(v, []string{"MYSQL_USERNAME", "DB_USERNAME"}, "devlop")
	EnvConfig.MYSQL_PASSWORD = getConfigValue(v, []string{"MYSQL_PASSWORD", "DB_PASSWORD"}, "123456")
	EnvConfig.MYSQL_DB_NAME = getConfigValue(v, []string{"MYSQL_DB_NAME", "DB_DATABASE"}, "devlop")

	// Redis配置
	EnvConfig.REDIS_HOST = v.GetString("REDIS_HOST")
	EnvConfig.REDIS_PASSWORD = v.GetString("REDIS_PASSWORD")
	EnvConfig.REDIS_PORT = v.GetString("REDIS_PORT")

	// DooTask配置
	EnvConfig.DOOTASK_DIR = v.GetString("DOOTASK_DIR")
	EnvConfig.DOOTASK_URL = v.GetString("DOOTASK_URL")

	// 第三方服务配置
	EnvConfig.YoudaoAppKey = v.GetString("YoudaoAppKey")
	EnvConfig.YoudaoAppSecret = v.GetString("YoudaoAppSecret")

	// 调试输出
	if EnvConfig.IsDev() {
		for _, key := range v.AllKeys() {
			fmt.Printf("读取配置项[ %s ] 值: %v\n", key, v.Get(key))
		}
	}
}
