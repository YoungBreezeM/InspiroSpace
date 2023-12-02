package api

import (
	"easygin/internal/conf"
	"easygin/internal/mid"
	"easygin/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Api struct {
	service  *service.Service
	conf     *conf.Config
	mid      *mid.Mid
	validate *validator.Validate
}

func InitApi(s *service.Service, e *gin.Engine, conf *conf.Config, m *mid.Mid) {

	a := Api{
		s,
		conf,
		m,
		validator.New(),
	}
	//register route
	a.RegisterWechatRoute(e)
	a.RegisterChatGTPRoute(e)
	a.RegisterUserRoute(e)
}
