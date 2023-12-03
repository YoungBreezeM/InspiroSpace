package api

import (
	"context"
	"easygin/internal/log"
	"easygin/internal/models"
	cgg "easygin/pkg/chatgtp_grpc"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ()

func (a *Api) RegisterChatGTPRoute(e *gin.Engine) {
	rg := e.Group("/chatgtp")

	rg.GET("/event/:chat_id", a.NotifyChat)

	rg.POST("/v1/chat/completions", a.mid.Auth, a.SubChat)

	rg.POST("/v1/chat/completions/test", a.SubChatTest)

}

func (a *Api) SubChat(c *gin.Context) {
	req := models.ApiChatRequest{}
	var openId string
	if o, ok := c.Get("openId"); ok {
		openId = o.(string)
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(200, models.R[string]{
			Status:  500,
			Data:    "",
			Message: err.Error(),
		})
		return
	}

	//参数校验
	if err := a.validate.Struct(&req); err != nil {
		c.JSON(200, models.R[string]{
			Status:  403,
			Data:    "",
			Message: err.Error(),
		})
		return
	}

	token, err := a.service.GetToken(c)
	if err != nil {
		c.JSON(200, models.R[string]{
			Status:  503,
			Data:    "",
			Message: "服务器资源紧张，请稍后再试！",
		})
		return
	}
	req.Token = token

	if err = a.service.SubChat(c, &req); err != nil {
		c.JSON(200, models.R[string]{
			Status:  500,
			Data:    "",
			Message: err.Error(),
		})
		return
	}

	go func() {
		if err := a.service.StoreChat(openId, req.Prompt); err != nil {
			log.Error(err)
		}
	}()

	c.JSON(200, models.R[string]{
		Status:  200,
		Data:    "",
		Message: "success",
	})

}

func (a *Api) SubChatTest(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 1)
		sendNotifyMsg(c, "hello")
	}
}

func sendNotifyMsg(c *gin.Context, msg string) {
	if _, err := c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", models.NewChatRespErr(msg))); err != nil {
		log.Error("write to client", err)
	}
	if _, err := c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", "[DONE]")); err != nil {
		log.Error("write to client", err)
	}
	c.Writer.Flush()
}

func (a *Api) NotifyChat(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	//
	chatId := c.Param("chat_id")
	if len(chatId) <= 0 {
		sendNotifyMsg(c, "chatId 不存在")
		return
	}

	req, err := a.service.NotifyChat(c, chatId)
	if err != nil {
		sendNotifyMsg(c, err.Error())
		return
	}

	log.Info(req)

	chatReq := cgg.ChatRequest{
		ChatId:          req.ChatId,
		Token:           req.Token,
		ConversationId:  req.ConversationId,
		ParentMessageId: req.ParentMessageId,
	}

	chatReq.Messages = append(chatReq.Messages, &cgg.ChatRequestMessage{
		Role:    "user",
		Content: req.Prompt,
	})
	//
	conn, err := grpc.Dial(a.service.ChatGtpConfig.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(fmt.Sprintf("did not connect %v", err))
		return
	}
	defer conn.Close()

	client := cgg.NewChatGTPServiceClient(conn)
	// Contact the server and print out its response.

	r, err := client.Chat(context.Background(), &chatReq)
	if err != nil {
		log.Error("chat func", err)
		return
	}

	for {
		select {
		case <-c.Writer.CloseNotify():
			log.Info(fmt.Sprintf("ChatId:%s is closed", chatId))
			return
		default:
			data, err := r.Recv()
			if err != nil && err != io.EOF {
				log.Info(fmt.Sprintf("ChatId:%s is closed", chatId))
				return
			}

			if data.Message == io.EOF.Error() {
				return
			}

			if len(data.Message) > 0 {
				log.Error("chatgtp error ", data.Message)
				if _, err = c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", models.NewChatRespErr(data.Message))); err != nil {
					log.Error("write to client", err)
				}
				if _, err = c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", "[DONE]")); err != nil {
					log.Error("write to client", err)
				}
				c.Writer.Flush()
				return
			}

			if _, err = c.Writer.WriteString(fmt.Sprintf("%s\n\n", data.Data)); err != nil {
				log.Error("write to client", err)
			}
			c.Writer.Flush()
		}

	}
}
