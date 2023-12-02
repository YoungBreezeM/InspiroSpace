package service

import (
	"easygin/internal/conf"
	"easygin/internal/dao"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
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
	s.SetToken("hello000")

	for i := 0; i < 100; i++ {
		go func(index int) {
			token, err := s.GetToken()
			if err != nil {
				log.Printf("%d error:%+v\n", index, err)
				return
			}
			fmt.Printf("token %d %s\n", index, token)
		}(i)
	}

	time.Sleep(time.Second * 2)

}
