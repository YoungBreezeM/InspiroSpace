package service

import (
	"context"
	"easygin/internal/conf"
	"easygin/internal/dao"
	"easygin/internal/models"
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	r, err := dao.InitRedis(&conf.RedisConfig{
		Host: "127.0.0.1",
		Port: 6379,
		DB:   0,
	})
	if err != nil {
		panic(err)
	}

	d := dao.Dao{
		Redis: r,
	}
	s := Service{Dao: &d}
	m := models.ApiChatRequest{
		ChatId: "da",
	}
	if err := s.SubChat(context.Background(), &m); err != nil {
		fmt.Println("sdsa", err)
	}

}
