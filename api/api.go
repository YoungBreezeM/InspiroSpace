package api

import (
	"easygin/internal/conf"
	"easygin/internal/mid"
	"easygin/internal/service"

	"github.com/gin-gonic/gin"
)

type Api struct {
	S    *service.Service
	conf *conf.Config
	M    *mid.Mid
	cq   *ApiChatRequest
}

func InitApi(s *service.Service, e *gin.Engine, conf *conf.Config, m *mid.Mid) {

	a := Api{s, conf, m, nil}
	//register route
	a.RegisterWechatRoute(e)
	a.RegisterChatGTPRoute(e)
	a.RegisterHelloRoute(e)
	a.RegisterUserRoute(e)
}
