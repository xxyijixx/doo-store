package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
)

var System SystemConfig

var (
	Language = []string{language.Chinese.String(), language.English.String()}
)

type SystemConfig struct {
	Port      int    `json:"port"`
	LogLevel  string `json:"log_level"`
	RemoteURL string `json:"remote_url"`
}

var defaultSystemConfig = SystemConfig{
	Port:      8080,
	LogLevel:  "info",
	RemoteURL: "http://localhost:9090",
}

func loadConfig(filename string) (SystemConfig, error) {
	var config SystemConfig

	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		config = defaultSystemConfig
		err = saveConfig(filename, config)
		if err != nil {
			return config, err
		}
	} else if err != nil {
		return config, err
	} else {
		err = json.Unmarshal(data, &config)
		if err != nil {
			return config, err
		}
	}

	return config, nil
}

func saveConfig(filename string, config SystemConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func systemInit() {
	configDir := "./conf"
	configFile := filepath.Join(configDir, "config.json")

	// 检查并创建 conf 目录（如果不存在）
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating config directory: %v\n", err)
			return
		}
	}

	System, err := loadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	fmt.Printf("Loaded config: %+v\n", System)
}
