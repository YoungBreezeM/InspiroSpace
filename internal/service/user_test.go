package service

import (
	"easygin/internal/conf"
	"easygin/internal/dao"
	"fmt"
	"testing"
)

func TestCreateToken(t *testing.T) {
	r, err := dao.InitRedis(&conf.RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "root",
		DB:       0,
	})
	if err != nil {
		panic(err)
	}

	d := dao.Dao{
		Redis: r,
	}
	s := Service{Dao: &d}
	s2 := s.CreateUser("hello")
	fmt.Println(s2)
}
