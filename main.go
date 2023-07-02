package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/eryajf/chatgpt-wecom/pkg/logger"
	"github.com/eryajf/chatgpt-wecom/pkg/process"
	"github.com/eryajf/chatgpt-wecom/public"
	"github.com/gin-gonic/gin"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
	"github.com/xen0n/go-workwx"
)

var wxclient *workwx.WorkwxApp

func init() {
	// 初始化加载配置，数据库，模板等
	public.InitSvc()
	// 指定日志等级
	logger.InitLogger(public.Config.LogLevel)
	var wx = workwx.New(public.Config.CorpId)

	wxclient = wx.WithApp(public.Config.AgentSecret, public.Config.AgentId)
	wxclient.SpawnAccessTokenRefresher() // 自动刷新token
}

func main() {
	StartHttp()
}

type dummyRxMessageHandler struct{}

var _ workwx.RxMessageHandler = dummyRxMessageHandler{}

// OnIncomingMessage 一条消息到来时的回调。
func (dummyRxMessageHandler) OnIncomingMessage(msg *workwx.RxMessage) error {
	go DoRequest(msg)
	return nil
}

func StartHttp() {
	app := gin.Default()
	app.GET("/*path", func(c *gin.Context) {
		// 获取到请求参数
		msgSignature := c.Query("msg_signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		echostr := c.Query("echostr")

		// 调用企业微信官方提供的接口进行解析校验
		wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(public.Config.ReceiveMsgToken, public.Config.ReceiveMsgKey, public.Config.CorpId, wxbizmsgcrypt.XmlType)
		echoStr, cryptErr := wxcpt.VerifyURL(msgSignature, timestamp, nonce, echostr)
		if nil != cryptErr {
			fmt.Println("verifyUrl fail", cryptErr)
		}
		c.String(200, string(echoStr))
	})
	hh, err := workwx.NewHTTPHandler(public.Config.ReceiveMsgToken, public.Config.ReceiveMsgKey, dummyRxMessageHandler{})
	if err != nil {
		logger.Fatal("init http handler failed: ", err)
	}

	app.POST("/*path", gin.WrapH(hh))
	port := ":" + public.Config.Port
	srv := &http.Server{
		Addr:    port,
		Handler: app,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		logger.Info("🚀 The HTTP Server is running on", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutting down server...")

	// 5秒后强制退出
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}
	logger.Info("Server exiting!")
}

func DoRequest(msgObj *workwx.RxMessage) {
	var msgInfo string
	resp := &workwx.Recipient{UserIDs: []string{msgObj.FromUserID}}
	if content, ok := msgObj.Text(); ok {
		msgInfo = content.GetContent()
	}
	// 再校验回调参数是否有价值
	if msgInfo == "" {
		logger.Warning("回调过来的内容为空，根据过往的经验，或许重新创建一下机器人，能解决这个问题")
		return
	}
	// 去除问题的前后空格
	msgInfo = strings.TrimSpace(msgInfo)
	if public.JudgeSensitiveWord(msgInfo) {
		logger.Info(fmt.Sprintf("🙋 %s提问的问题中包含敏感词汇，userid：%#v，消息: %#v", msgObj.FromUserID, msgObj.FromUserID, msgInfo))
		return
	}
	// 打印企微回调过来的请求明细，调试时打开
	logger.Debug(fmt.Sprintf("wecom callback parameters: %#v", msgObj))

	if !public.JudgeUsers(msgObj.FromUserID) && !public.JudgeAdminUsers(msgObj.FromUserID) && msgObj.FromUserID != "" {
		logger.Info(fmt.Sprintf("🙋 %s身份信息未被验证通过，userid：%#v，消息: %#v", msgObj.FromUserID, msgObj.FromUserID, msgInfo))
		return
	}
	if len(msgInfo) == 0 || msgInfo == "帮助" {

		// 欢迎信息
		err := wxclient.SendMarkdownMessage(resp, public.Config.Help, false)
		if err != nil {
			fmt.Printf("send err:%v", err)
		}
	} else {
		logger.Info(fmt.Sprintf("🙋 %s发起的问题: %#v", msgObj.FromUserID, msgInfo))
		// 除去帮助之外的逻辑分流在这里处理
		switch {
		default:
			var err error
			msgInfo, err = process.GeneratePrompt(msgInfo)
			// err不为空：提示词之后没有文本 -> 直接返回提示词所代表的内容
			if err != nil {
				_ = wxclient.SendMarkdownMessage(resp, msgInfo, false)
				return
			}
			err = process.ProcessRequest(msgObj, wxclient)
			if err != nil {
				logger.Warning(fmt.Errorf("process request: %v", err))
				return
			}
			return
		}
	}
}
