package compose

import (
	"doo-store/backend/config"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var envMap = map[string]string{
	"DOOTASK_NETWORK_NAME": config.EnvConfig.App().NETWORK_NAME,
}

func ReplaceEnvVariables(input string) string {
	// 使用正则表达式匹配 ${xxx} 格式的环境变量
	re := regexp.MustCompile(`\${([^}]+)}`)

	// 替换所有匹配的环境变量
	return re.ReplaceAllStringFunc(input, func(match string) string {
		// 提取变量名
		varName := strings.Trim(match, "${}")
		// 从 map 中查找该变量的值
		if envValue, exists := envMap[varName]; exists {
			return envValue
		}
		// 如果没有找到变量值，返回原占位符
		return match
	})
}

// parseEnvContent 将 envContent 解析为键值对
func parseEnvContent(envContent string) (map[string]string, error) {
	envMap := make(map[string]string)

	// 按行分割 envContent
	lines := strings.Split(envContent, "\n")
	for _, line := range lines {
		// 忽略空行
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 按等号分割键值对
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid envContent format: expected K=V")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 检查键是否为空
		if key == "" {
			return nil, errors.New("empty key in envContent")
		}

		// 存入 map
		envMap[key] = value
	}

	return envMap, nil
}

// replaceEnvVars 替换 content 中的环境变量
func replaceEnvVars(content string, envMap map[string]string) string {
	// 遍历 envMap，替换 content 中的 ${VAR}
	for key, value := range envMap {
		placeholder := fmt.Sprintf("${%s}", key)
		content = strings.ReplaceAll(content, placeholder, value)
	}

	return content
}
