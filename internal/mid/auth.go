package mid

import (
	"easygin/internal/models"
	"easygin/internal/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Token struct {
	OpenId string
	Key    string
}

func (m *Mid) Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) <= 0 {
		c.JSON(200, models.R[string]{
			Status:  401,
			Data:    "",
			Message: "认证失败",
		})
		c.Abort()
		return
	}

	//解析token
	fmt.Println(token)
	token = strings.Replace(token, " ", "+", 10)
	t := Token{}
	if err := utils.ParaseToken[Token](token, m.S.Config.AppConfig.TokenKey, &t); err != nil {
		c.JSON(200, models.R[string]{
			Status:  401,
			Data:    "",
			Message: "解析token失败",
		})
		c.Abort()
		return
	}

	s, err := m.S.Redis.Get(c, t.OpenId).Result()
	if err != nil {
		c.JSON(200, models.R[string]{
			Status:  401,
			Data:    "",
			Message: "链接已过期,请重新获取链接!",
		})
		c.Abort()
		return
	}

	if t.Key != s {
		c.JSON(200, models.R[string]{
			Status:  401,
			Data:    "",
			Message: "key 无效",
		})
		c.Abort()
		return
	}

	c.Set("openId", t.OpenId)

	c.Next()
}
