package process

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/eryajf/chatgpt-wecom/pkg/db"
	"github.com/eryajf/chatgpt-wecom/pkg/logger"
	"github.com/eryajf/chatgpt-wecom/public"
	"github.com/solywsh/chatgpt"
	"github.com/xen0n/go-workwx"
)

// ProcessRequest 分析处理请求逻辑
func ProcessRequest(rmsg *workwx.RxMessage, wxclient *workwx.WorkwxApp) error {
	var msgInfo string
	resp := &workwx.Recipient{UserIDs: []string{rmsg.FromUserID}}
	if content, ok := rmsg.Text(); ok {
		msgInfo = content.GetContent()
	}
	if CheckRequestTimes(rmsg, wxclient) {
		content := strings.TrimSpace(msgInfo)
		timeoutStr := ""
		if content != public.Config.DefaultMode {
			timeoutStr = fmt.Sprintf("\n\n>%s 后将恢复默认聊天模式：%s", FormatTimeDuation(public.Config.SessionTimeout), public.Config.DefaultMode)
		}
		switch content {
		case "单聊":
			public.UserService.SetUserMode(rmsg.FromUserID, content)
			err := wxclient.SendMarkdownMessage(resp, fmt.Sprintf("**[Concentrate] 现在进入与 %s 的单聊模式**%s", rmsg.FromUserID, timeoutStr), false)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "串聊":
			public.UserService.SetUserMode(rmsg.FromUserID, content)
			err := wxclient.SendMarkdownMessage(resp, fmt.Sprintf("**[Concentrate] 现在进入与 %s 的串聊模式**%s", rmsg.FromUserID, timeoutStr), false)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "重置", "退出", "结束":
			// 重置用户对话模式
			public.UserService.ClearUserMode(rmsg.FromUserID)
			// 清空用户对话上下文
			public.UserService.ClearUserSessionContext(rmsg.FromUserID)
			// 清空用户对话的答案ID
			public.UserService.ClearAnswerID(rmsg.FromUserID, rmsg.FromUserID)
			err := wxclient.SendMarkdownMessage(resp, fmt.Sprintf("[RecyclingSymbol]已重置与**%s** 的对话模式\n\n> 可以开始新的对话 [Bubble]", rmsg.FromUserID), false)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		case "模板":
			var title string
			for _, v := range *public.Prompt {
				title = title + v.Title + " | "
			}
			err := wxclient.SendMarkdownMessage(resp, fmt.Sprintf("%s 您好，当前程序内置集成了这些提示词：\n\n-----\n\n| %s \n\n-----\n\n您可以选择某个提示词作为对话内容的开头。\n\n以周报为例，可发送\"#周报 我本周用Go写了一个企微集成ChatGPT的聊天应用\"，可将工作内容填充为一篇完整的周报。\n\n-----\n\n若您不清楚某个提示词的所代表的含义，您可以直接发送提示词，例如直接发送\"#周报\"", rmsg.FromUserID, title), false)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
		default:
			if public.FirstCheck(rmsg.FromUserID) {
				return Do("串聊", rmsg, wxclient)
			} else {
				return Do("单聊", rmsg, wxclient)
			}
		}
	}
	return nil
}

// 执行处理请求
func Do(mode string, rmsg *workwx.RxMessage, wxclient *workwx.WorkwxApp) error {
	var msgInfo string
	resp := &workwx.Recipient{UserIDs: []string{rmsg.FromUserID}}
	if content, ok := rmsg.Text(); ok {
		msgInfo = content.GetContent()
	}
	// 先把模式注入
	public.UserService.SetUserMode(rmsg.FromUserID, mode)
	switch mode {
	case "单聊":
		qObj := db.Chat{
			Username:      rmsg.FromUserID,
			Source:        rmsg.FromUserID,
			ChatType:      db.Q,
			ParentContent: 0,
			Content:       msgInfo,
		}
		qid, err := qObj.Add()
		if err != nil {
			logger.Error("往MySQL新增数据失败,错误信息：", err)
		}
		reply, err := chatgpt.SingleQa(msgInfo, rmsg.FromUserID)
		if err != nil {
			logger.Info(fmt.Errorf("gpt request error: %v", err))
			if strings.Contains(fmt.Sprintf("%v", err), "maximum question length exceeded") {
				public.UserService.ClearUserSessionContext(rmsg.FromUserID)
				err := wxclient.SendMarkdownMessage(resp, fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v\n\n> 已超过最大文本限制，请缩短提问文字的字数。", err), false)
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			} else {
				err := wxclient.SendMarkdownMessage(resp, fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v", err), false)
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			}
		}
		if reply == "" {
			logger.Warning(fmt.Errorf("get gpt result falied: %v", err))
			return nil
		} else {
			reply = strings.TrimSpace(reply)
			reply = strings.Trim(reply, "\n")
			aObj := db.Chat{
				Username:      rmsg.FromUserID,
				Source:        rmsg.FromUserID,
				ChatType:      db.A,
				ParentContent: qid,
				Content:       reply,
			}
			_, err := aObj.Add()
			if err != nil {
				logger.Error("往MySQL新增数据失败,错误信息：", err)
			}
			logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.FromUserID, reply))
			if public.JudgeSensitiveWord(reply) {
				reply = public.SolveSensitiveWord(reply)
			}
			// 回复@我的用户
			err = wxclient.SendMarkdownMessage(resp, FormatMarkdown(reply), false)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return err
			}
		}
	case "串聊":
		lastAid := public.UserService.GetAnswerID(rmsg.FromUserID, rmsg.FromUserID)
		qObj := db.Chat{
			Username:      rmsg.FromUserID,
			Source:        rmsg.FromUserID,
			ChatType:      db.Q,
			ParentContent: lastAid,
			Content:       msgInfo,
		}
		qid, err := qObj.Add()
		if err != nil {
			logger.Error("往MySQL新增数据失败,错误信息：", err)
		}
		cli, reply, err := chatgpt.ContextQa(msgInfo, rmsg.FromUserID)
		if err != nil {
			logger.Info(fmt.Sprintf("gpt request error: %v", err))
			if strings.Contains(fmt.Sprintf("%v", err), "maximum text length exceeded") {
				public.UserService.ClearUserSessionContext(rmsg.FromUserID)
				err = wxclient.SendMarkdownMessage(resp, fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v\n\n> 串聊已超过最大文本限制，对话已重置，请重新发起。", err), false)
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			} else {
				err = wxclient.SendMarkdownMessage(resp, fmt.Sprintf("[Wrong] 请求 OpenAI 失败了\n\n> 错误信息:%v", err), false)
				if err != nil {
					logger.Warning(fmt.Errorf("send message error: %v", err))
					return err
				}
			}
		}
		if reply == "" {
			logger.Warning(fmt.Errorf("get gpt result falied: %v", err))
			return nil
		} else {
			reply = strings.TrimSpace(reply)
			reply = strings.Trim(reply, "\n")
			aObj := db.Chat{
				Username:      rmsg.FromUserID,
				Source:        rmsg.FromUserID,
				ChatType:      db.A,
				ParentContent: qid,
				Content:       reply,
			}
			aid, err := aObj.Add()
			if err != nil {
				logger.Error("往MySQL新增数据失败,错误信息：", err)
			}
			// 将当前回答的ID放入缓存
			public.UserService.SetAnswerID(rmsg.FromUserID, string(rmsg.MsgType), aid)
			logger.Info(fmt.Sprintf("🤖 %s得到的答案: %#v", rmsg.FromUserID, reply))
			if public.JudgeSensitiveWord(reply) {
				reply = public.SolveSensitiveWord(reply)
			}
			// 回复@我的用户
			err = wxclient.SendMarkdownMessage(resp, FormatMarkdown(reply), false)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return err
			}
			_ = cli.ChatContext.SaveConversation(rmsg.FromUserID)
		}
	default:

	}
	return nil
}

// FormatTimeDuation 格式化时间
// 主要提示单聊/群聊切换时多久后恢复默认聊天模式
func FormatTimeDuation(duration time.Duration) string {
	minutes := int64(duration.Minutes())
	seconds := int64(duration.Seconds()) - minutes*60
	var timeoutStr string
	if seconds == 0 {
		timeoutStr = fmt.Sprintf("%d分钟", minutes)
	} else {
		timeoutStr = fmt.Sprintf("%d分%d秒", minutes, seconds)
	}
	return timeoutStr
}

// FormatMarkdown 格式化Markdown
// 主要修复ChatGPT返回多行代码块，企微会将代码块中的#当作Markdown语法里的标题来处理，进行转义；如果Markdown格式内存在html，将Markdown中的html标签转义
// 代码块缩进问题暂无法解决，因不管是四个空格，还是Tab，在企微上均会顶格显示，建议复制代码后用IDE进行代码格式化，针对缩进严格的语言，例如Python，不确定的建议手机端查看下代码块的缩进
func FormatMarkdown(md string) string {
	lines := strings.Split(md, "\n")
	codeblock := false
	existHtml := strings.Contains(md, "<")

	for i, line := range lines {
		if strings.HasPrefix(line, "```") {
			codeblock = !codeblock
		}
		if codeblock {
			lines[i] = strings.ReplaceAll(line, "#", "\\#")
		} else if existHtml {
			lines[i] = html.EscapeString(line)
		}
	}

	return strings.Join(lines, "\n")
}

// CheckRequestTimes 分析处理请求逻辑
// 主要提供单日请求限额的功能
func CheckRequestTimes(rmsg *workwx.RxMessage, wxclient *workwx.WorkwxApp) bool {
	resp := &workwx.Recipient{UserIDs: []string{rmsg.FromUserID}}
	if public.Config.MaxRequest == 0 {
		return true
	}
	count := public.UserService.GetUseRequestCount(rmsg.FromUserID)
	// 用户是管理员或VIP用户，不判断访问次数是否超过限制
	if public.JudgeAdminUsers(rmsg.FromUserID) || public.JudgeVipUsers(rmsg.FromUserID) {
		return true
	} else {
		// 用户不是管理员和VIP用户，判断访问次数是否超过限制
		if count >= public.Config.MaxRequest {
			logger.Info(fmt.Sprintf("亲爱的: %s，您今日请求次数已达上限，请明天再来，交互发问资源有限，请务必斟酌您的问题，给您带来不便，敬请谅解!", rmsg.FromUserID))
			err := wxclient.SendMarkdownMessage(resp, fmt.Sprintf("[Staple] **一个好的问题，胜过十个好的答案！** \n\n亲爱的%s:\n\n您今日请求次数已达上限，请明天再来，交互发问资源有限，请务必斟酌您的问题，给您带来不便，敬请谅解！\n\n如有需要，可联系管理员升级为VIP用户。", rmsg.FromUserID), false)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
			}
			return false
		}
	}
	// 访问次数未超过限制，将计数加1
	public.UserService.SetUseRequestCount(rmsg.FromUserID, count+1)
	return true
}
