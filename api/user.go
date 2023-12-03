package api

import (
	"easygin/internal/mid"
	"easygin/internal/models"
	"easygin/internal/utils"
	"os"

	"github.com/gin-gonic/gin"
)

func (a *Api) RegisterUserRoute(e *gin.Engine) {
	rg := e.Group("/user")

	rg.GET("/token", a.GetToken)
}

func (a *Api) GetToken(c *gin.Context) {
	s := os.Getenv("GIN_MODE")

	if s == "release" {
		c.JSON(200, "fail")
		return
	}

	token := mid.Token{
		OpenId: "test",
		Key:    a.service.CreateUser(c, "test"),
	}

	tokenStr, err := utils.CreateToken(token, a.conf.AppConfig.TokenKey)
	if err != nil {
		c.JSON(200, models.R[string]{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, models.R[string]{
		Status:  0,
		Message: "successful",
		Data:    tokenStr,
	})
}
