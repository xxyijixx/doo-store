package docker

import (
	"errors"
	"regexp"
)

func ParseError(stdoutStr string, err error) (string, error) {
	// 使用正则表达式匹配 "Error response from daemon" 后的内容
	re := regexp.MustCompile(`Error response from daemon: (.*)`)
	match := re.FindStringSubmatch(stdoutStr)

	errMsg := ""
	if len(match) > 1 {
		// 如果匹配成功，取匹配到的内容
		errMsg = match[1]
	} else if len(stdoutStr) > 0 {
		// 如果没有匹配到，返回完整的 stderr 内容
		errMsg = stdoutStr
	}

	return "", errors.New(errMsg)
}
