package public

import (
	"strings"
)

func FirstCheck(userid string) bool {
	lc := UserService.GetUserMode(userid)
	if lc == "" {
		if Config.DefaultMode == "串聊" {
			return true
		} else {
			return false
		}
	}
	if lc != "" && strings.Contains(lc, "串聊") {
		return true
	}
	return false
}
