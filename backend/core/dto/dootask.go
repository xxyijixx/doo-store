package dto

import (
	"fmt"
	"strconv"
	"strings"
)

// UserInfoResp 返回用户信息
type UserInfoResp struct {
	*UserBasicResp
	Token            string   `json:"token"`             // token
	Identity         []string `json:"identity"`          // 身份
	Tel              string   `json:"tel"`               // 手机号
	Changepass       int      `json:"changepass"`        // 是否需要修改密码 1-是 0-否
	CreatedIp        string   `json:"created_ip"`        // 创建ip
	EmailVerity      int      `json:"email_verity"`      // 邮箱是否验证 1-是 0-否
	LastAt           string   `json:"last_at"`           // 最后登录时间
	LastIp           string   `json:"last_ip"`           // 最后登录ip
	LineIp           string   `json:"line_ip"`           // 在线ip
	LoginNum         int      `json:"login_num"`         // 登录次数
	NicknameOriginal string   `json:"nickname_original"` // 昵称原始值
	TaskDialogId     int      `json:"task_dialog_id"`    // 任务对话框id
	CreatedAt        string   `json:"created_at"`        // 创建时间
	DepartmentOwner  bool     `json:"department_owner"`  // 是否是部门负责人 false-不是 true-是
	OkrAdminOwner    bool     `json:"okr_admin_owner"`   // okr普通人员是否拥有管理员有权限 false-不是 true-是
}

// 用户基础信息
type UserBasicResp struct {
	Az             string `json:"az"`              // 首字母
	Bot            int    `json:"bot"`             // 是否是机器人 1-是 0-否
	Department     []int  `json:"department"`      // 部门
	DepartmentName string `json:"department_name"` // 部门名称
	DisableAt      string `json:"disable_at"`      // 禁用时间
	Userid         int    `json:"userid"`          // 用户id
	Nickname       string `json:"nickname"`        // 昵称
	Userimg        string `json:"userimg"`         // 头像
	Email          string `json:"email"`           // 邮箱
	LineAt         string `json:"line_at"`         // 在线时间
	Pinyin         string `json:"pinyin"`          // 拼音
	Profession     string `json:"profession"`      // 职业
}

// 判断是否是管理员
func (u *UserInfoResp) IsAdmin() bool {
	for _, v := range u.Identity {
		if v == "admin" {
			return true
		}
	}
	return false
}

type VersionInfoResp struct {
	Version string `json:"version"`
	Publish struct {
		Provider string `json:"provider"`
		URL      string `json:"url"`
	} `json:"publish"`
}

func (v *VersionInfoResp) CheckVersion(requiredVersion string) (bool, error) {
	current, err := v.parseVersion(v.Version)
	if err != nil {
		return false, err
	}

	required, err := v.parseVersion(requiredVersion)
	if err != nil {
		return false, err
	}

	// 比较每一个版本位
	for i := 0; i < 3; i++ {
		if current[i] > required[i] {
			return true, nil
		} else if current[i] < required[i] {
			return false, nil
		}
	}

	return true, nil
}

// parseVersion 解析版本字符串，返回主版本号、次版本号、补丁版本号
func (v *VersionInfoResp) parseVersion(version string) ([]int, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("版本号格式不正确: %s", version)
	}

	var numbers []int
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("解析版本号失败: %s", version)
		}
		numbers = append(numbers, num)
	}

	return numbers, nil
}
