package service

import (
	"easygin/internal/conf"
	"easygin/internal/dao"
)

type Service struct {
	*dao.Dao
	*conf.Config
}

func InitService(d *dao.Dao, c *conf.Config) *Service {
	return &Service{d, c}
}
