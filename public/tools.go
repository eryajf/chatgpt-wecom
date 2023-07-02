package public

import (
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

// 将内容写入到文件，如果文件名带路径，则会判断路径是否存在，不存在则创建
func WriteToFile(path string, data []byte) error {
	tmp := strings.Split(path, "/")
	if len(tmp) > 0 {
		tmp = tmp[:len(tmp)-1]
	}

	err := os.MkdirAll(strings.Join(tmp, "/"), os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0755)
	if err != nil {
		return err
	}
	return nil
}

// JudgeUsers 判断用户是否在白名单
func JudgeUsers(s string) bool {
	// 优先判断黑名单，黑名单用户返回：不在白名单
	if len(Config.DenyUsers) != 0 {
		for _, v := range Config.DenyUsers {
			if v == s {
				return false
			}
		}
	}
	// 白名单配置逻辑处理
	if len(Config.AllowUsers) == 0 {
		return true
	}
	for _, v := range Config.AllowUsers {
		if v == s {
			return true
		}
	}
	return false
}

// JudgeAdminUsers 判断用户是否为系统管理员
func JudgeAdminUsers(s string) bool {
	// 如果secret或者用户的userid都为空的话，那么默认没有管理员
	if s == "" {
		return false
	}
	// 如果没有指定，则没有人是管理员
	if len(Config.AdminUsers) == 0 {
		return false
	}
	for _, v := range Config.AdminUsers {
		if v == s {
			return true
		}
	}
	return false
}

// JudgeVipUsers 判断用户是否为VIP用户
func JudgeVipUsers(s string) bool {
	// 如果secret或者用户的userid都为空的话，那么默认不是VIP用户
	if s == "" {
		return false
	}
	// 管理员默认是VIP用户
	for _, v := range Config.AdminUsers {
		if v == s {
			return true
		}
	}
	// 如果没有指定，则没有人是VIP用户
	if len(Config.VipUsers) == 0 {
		return false
	}
	for _, v := range Config.VipUsers {
		if v == s {
			return true
		}
	}
	return false
}

func GetReadTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// JudgeSensitiveWord 判断内容是否包含敏感词
func JudgeSensitiveWord(s string) bool {
	if len(Config.SensitiveWords) == 0 {
		return false
	}
	for _, v := range Config.SensitiveWords {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

// SolveSensitiveWord 将敏感词用 🚫 占位
func SolveSensitiveWord(s string) string {
	for _, v := range Config.SensitiveWords {
		if strings.Contains(s, v) {
			return strings.Replace(s, v, printStars(utf8.RuneCountInString(v)), -1)
		}
	}
	return s
}

// 将对应敏感词替换为 🚫
func printStars(num int) string {
	s := ""
	for i := 0; i < num; i++ {
		s += "🚫"
	}
	return s
}
