package common

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// RandString 生成随机字符串
func RandString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RandNum 生成随机数
func RandNum(min, max int) int {
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// 字符串转int
func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// 检查字符串s是否以prefix开头
func LeftExists(s string, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// 删除字符串s中所有前缀为prefix的子串
func LeftDelete(s string, prefix string) string {
	for LeftExists(s, prefix) {
		s = s[len(prefix):]
	}
	return s
}

// 计算两个字符串切片之间的差集
func StringDiff(slice1, slice2 []string) []string {
	diff := []string{}
	for _, s1 := range slice1 {
		if !StringContains(slice2, s1) {
			diff = append(diff, s1)
		}
	}
	return diff
}

// 检查字符串s是否存在于slice中
func StringContains(slice []string, s string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}
	return false
}

// 检查字符串是否是合法的 MAC 地址格式（例如：00:11:22:33:44:55）
func IsMAC(mac string) bool {
	regex := regexp.MustCompile("^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$")
	return regex.MatchString(mac)
}

// 打散字符串，只留为数字的项
func ExplodeInt(delimiter string, str interface{}, reInt bool) []int {
	var array []string
	if str == nil {
		str = delimiter
		delimiter = ","
	}
	// 判断类型
	switch v := str.(type) {
	case string:
		array = strings.Split(v, delimiter)
	case []string:
		array = v
	case []int:
		return v
	}
	return ArrayRetainInt(array, reInt)
}

// 数组只保留数字的
func ArrayRetainInt(arr []string, reInt bool) []int {
	var result []int
	for _, v := range arr {
		// 如果不是数字，则剔除
		if _, err := strconv.Atoi(v); err != nil {
			continue
		}
		// 如果需要格式化，则格式化
		if reInt {
			result = append(result, StringToInt(v))
		}
	}
	return result
}

// GetMiddle 截取指定字符串
func GetMiddle(str, ta, tb string) string {
	if ta != "" && strings.Contains(str, ta) {
		str = str[strings.Index(str, ta)+len(ta):]
	}
	if tb != "" && strings.Contains(str, tb) {
		str = str[:strings.Index(str, tb)]
	}
	return str
}

// IsChineseCharCountValid 检查字符串中的中文字符数量是否超过255个
func IsChineseCharCountValid(str string) bool {
	return utf8.RuneCountInString(str) <= 255
}
