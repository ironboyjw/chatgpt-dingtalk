package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/eryajf/chatgpt-dingtalk/pkg/dingbot"
	"github.com/eryajf/chatgpt-dingtalk/pkg/logger"
	"github.com/eryajf/chatgpt-dingtalk/pkg/process"
	"github.com/eryajf/chatgpt-dingtalk/public"
	"github.com/xgfone/ship/v5"
)

func init() {
	public.InitSvc()
}
func main() {
	Start()
}

<<<<<<< HEAD
=======
var Welcome string = `Commands:
=================================
🙋 单聊 👉 单独聊天
📣 串聊 👉 带上下文聊天
🔃 重置 👉 重置带上下文聊天
💵 余额 👉 查询剩余额度
🚀 帮助 👉 显示帮助信息
🌈 模板 👉 内置的prompt
🎨 图片 👉 根据prompt生成图片
=================================
🚜 例：@我发送 空 或 帮助 将返回此帮助信息
💪 Power By https://github.com/eryajf/chatgpt-dingtalk
`

>>>>>>> parent of 71a464b (perf: 当使用prompt但内容为空时，直接返回prompt的内容 (#138))
func Start() {
	app := ship.Default()
	app.Route("/").POST(func(c *ship.Context) error {
		var msgObj dingbot.ReceiveMsg
		err := c.Bind(&msgObj)
		if err != nil {
			return ship.ErrBadRequest.New(fmt.Errorf("bind to receivemsg failed : %v", err))
		}
		if msgObj.Text.Content == "" || msgObj.ChatbotUserID == "" {
			logger.Warning("从钉钉回调过来的内容为空，根据过往的经验，或许重新创建一下机器人，能解决这个问题")
			return ship.ErrBadRequest.New(fmt.Errorf("从钉钉回调过来的内容为空，根据过往的经验，或许重新创建一下机器人，能解决这个问题"))
		}

		// 打印钉钉回调过来的请求明细
		logger.Info(fmt.Sprintf("dingtalk callback parameters: %#v", msgObj))
		// TODO: 校验请求
		if len(msgObj.Text.Content) == 1 || strings.TrimSpace(msgObj.Text.Content) == "帮助" {
			// 欢迎信息
			_, err := msgObj.ReplyToDingtalk(string(dingbot.TEXT), Welcome)
			if err != nil {
				logger.Warning(fmt.Errorf("send message error: %v", err))
				return ship.ErrBadRequest.New(fmt.Errorf("send message error: %v", err))
			}
		} else {
			// 除去帮助之外的逻辑分流在这里处理
			switch {
			case strings.HasPrefix(strings.TrimSpace(msgObj.Text.Content), "#图片"):
				return process.ImageGenerate(&msgObj)
			default:
				msgObj.Text.Content = process.GeneratePrompt(strings.TrimSpace(msgObj.Text.Content))
				logger.Info(fmt.Sprintf("after generate prompt: %#v", msgObj.Text.Content))
				return process.ProcessRequest(&msgObj)
			}
		}
		return nil
	})
	// 解析生成后的图片
	app.Route("/images/:filename").GET(func(c *ship.Context) error {
		filename := c.Param("filename")
		root := "./images/"
		return c.File(filepath.Join(root, filename))
	})

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
	// 启动服务器
	ship.StartServer(":8090", app)
}

var Welcome string = `# 发送信息

若您想给机器人发送信息，请选择：

1. 在本机器人所在群里@机器人；
2. 点击机器人的头像后，再点击"发消息"。

机器人收到您的信息后，默认会交给chatgpt进行处理。除非，您发送的内容是7个**系统指令**之一。

-----

# 系统指令

系统指令是一些特殊的词语，当您向机器人发送这些词语时，会触发对应的功能：

**单聊**：每条消息都是单独的对话，不包含上下文

**串聊**：对话会携带上下文，除非您主动重置对话或对话长度超过限制

**重置**：重置上下文

**余额**：查询机器人所用OpenAI账号的余额

**模板**：查询机器人内置的快捷模板

**图片**：查看如何根据提示词生成图片

**帮助**：重新获取帮助信息

-----

# 友情提示

使用"串聊模式"会显著加快机器人所用账号的余额消耗速度。

因此，若无保留上下文的需求，建议使用"单聊模式"。

即使有保留上下文的需求，也应适时使用"重置"指令来重置上下文。

-----

# 项目地址

本项目已在GitHub开源，[查看源代码](https://github.com/eryajf/chatgpt-dingtalk)。
`
