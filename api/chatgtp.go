package api

import (
	"context"
	"easygin/api/pb"
	"easygin/internal/log"
	"easygin/internal/models"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	stopTimer = make(chan struct{})
)

type ApiChatRequest struct {
	ChatId          string `json:"chat_id"`
	Prompt          string `json:"prompt"`
	ParentMessageId string `json:"parent_message_id"`
	ConversationId  string `json:"conversation_id"`
	T               *time.Timer
	E               chan struct{}
	sync.Mutex
}

func (a *Api) RegisterChatGTPRoute(e *gin.Engine) {
	rg := e.Group("/chatgtp")

	rg.GET("/event/:chat_id", a.NotifyChat)

	rg.POST("/v1/chat/completions", a.M.Auth, a.SubChat)

}

func (a *Api) SubChat(c *gin.Context) {
	req := ApiChatRequest{}
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

	a.cq = &req
	if !a.cq.TryLock() {
		c.JSON(200, models.R[string]{
			Status:  503,
			Data:    "",
			Message: "服务器繁忙，请稍后再试！",
		})
		return
	}

	a.cq.T = time.NewTimer(10 * time.Second)
	a.cq.E = make(chan struct{})
	//occupancy clear
	go func(t *time.Timer, e chan struct{}) {
		defer func() {
			t.Stop()
		}()

		for {
			select {
			case <-t.C:
				log.Warnf("chat id %s is timeout", a.cq.ChatId)
				a.cq.Unlock()
				a.cq = nil
				return
			case <-e:
				log.Info("start read chatgtp content")
				return
			}

		}

	}(a.cq.T, a.cq.E)

	if err := a.S.StoreChat(openId, req.Prompt); err != nil {
		log.Error(err)
	}

	c.JSON(200, models.R[string]{
		Status:  200,
		Data:    "",
		Message: "success",
	})

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

	if a.cq == nil {
		sendNotifyMsg(c, "chatId 位注册")
		return
	}

	if chatId != a.cq.ChatId {
		sendNotifyMsg(c, "请求不存在")
		return
	}

	defer func() {
		if a.cq != nil {
			a.cq.Unlock()
			a.cq = nil
		}
	}()

	a.cq.E <- struct{}{}
	chatReq := pb.ChatRequest{
		ChatId:          chatId,
		Token:           a.conf.ChatGtpConfig.Token,
		ConversationId:  "",
		ParentMessageId: "",
	}

	if len(a.cq.ConversationId) > 0 {
		chatReq.ConversationId = a.cq.ConversationId
	}

	if len(a.cq.ParentMessageId) > 0 {
		chatReq.ParentMessageId = a.cq.ParentMessageId
	}

	chatReq.Messages = append(chatReq.Messages, &pb.ChatRequestMessage{
		Role:    "user",
		Content: a.cq.Prompt,
	})
	//
	conn, err := grpc.Dial(a.S.ChatGtpConfig.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(fmt.Sprintf("did not connect %v", err))
		return
	}
	defer conn.Close()

	client := pb.NewChatGTPServiceClient(conn)
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
