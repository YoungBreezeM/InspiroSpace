package api

import (
	"crypto/sha1"
	"easygin/internal/mid"
	"easygin/internal/models"
	"easygin/internal/utils"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CallbackRequest struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
	EventKey     string   `xml:"EventKey"`
	Ticket       string   `xml:"Ticket"`
	Content      string   `xml:"Content"`
}

type CallbackReply struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
}

func (a *Api) RegisterWechatRoute(e *gin.Engine) {
	rg := e.Group("/wechat")

	rg.GET("/notify", a.AuthToken)
	rg.POST("/notify", a.Notify)
}

func (a *Api) AuthToken(c *gin.Context) {
	queryValues := c.Request.URL.Query()

	signature := queryValues.Get("signature")
	timestamp := queryValues.Get("timestamp")
	nonce := queryValues.Get("nonce")
	echoStr := queryValues.Get("echostr")

	// 将 token、timestamp 和 nonce 放入一个切片中，并按照字典顺序排序
	params := []string{a.conf.WechatConfig.Token, timestamp, nonce}
	sort.Strings(params)

	// 将排序后的参数拼接成一个字符串
	str := strings.Join(params, "")

	// 对拼接后的字符串进行 SHA1 计算
	h := sha1.New()
	_, err := io.WriteString(h, str)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	sha1Sum := hex.EncodeToString(h.Sum(nil))

	// 将计算得到的 SHA1 值与 signature 进行比较
	if sha1Sum == signature {
		// Token 验证通过，返回 echostr 参数
		c.String(http.StatusOK, echoStr)
	} else {
		// Token 验证失败
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func (a *Api) Notify(c *gin.Context) {
	c.Header("Content-Type", "application/xml")
	callbackMsg := CallbackRequest{}
	var replyMsg CallbackReply

	if err := c.ShouldBindXML(&callbackMsg); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if callbackMsg.MsgType == "text" {

		if !a.conf.AppConfig.Enable {
			replyMsg = CallbackReply{
				ToUserName:   callbackMsg.FromUserName,
				FromUserName: callbackMsg.ToUserName,
				CreateTime:   time.Now().Unix(),
				MsgType:      "text",
				Content:      "系统正在维护中，请稍后再试，谢谢！",
			}
			c.XML(http.StatusOK, replyMsg)
			return
		}

		token := mid.Token{
			OpenId: callbackMsg.FromUserName,
			Key:    a.S.CreateUser(callbackMsg.FromUserName),
		}

		tokenStr, err := utils.CreateToken(token, a.conf.AppConfig.TokenKey)
		if err != nil {
			c.JSON(200, models.R[string]{
				Message: err.Error(),
			})
			return
		}

		//
		if callbackMsg.Content == "chatgtp" {
			replyMsg = CallbackReply{
				ToUserName:   callbackMsg.FromUserName,
				FromUserName: callbackMsg.ToUserName,
				CreateTime:   time.Now().Unix(),
				MsgType:      "text",
				Content:      fmt.Sprintf("点开下面这个链接进入网页\n%s/?token=%s", a.conf.WechatConfig.CallbackUrl, tokenStr),
			}
			//
			c.XML(http.StatusOK, replyMsg)
		} else {
			replyMsg = CallbackReply{
				ToUserName:   callbackMsg.FromUserName,
				FromUserName: callbackMsg.ToUserName,
				CreateTime:   time.Now().Unix(),
				MsgType:      "text",
				Content:      "您输入的关键词不对哦！目前可用关键词为[chatgtp]",
			}
			c.XML(http.StatusOK, replyMsg)
		}

	}
}
