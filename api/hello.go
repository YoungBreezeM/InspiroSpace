package api

import (
	"easygin/internal/log"

	"github.com/gin-gonic/gin"
)

func (a *Api) RegisterHelloRoute(e *gin.Engine) {
	rg := e.Group("/hello")

	rg.GET("/", a.SayHello)
}

func (a *Api) SayHello(c *gin.Context) {
	log.Debug("hello")
	c.JSON(200, "ok")
}
