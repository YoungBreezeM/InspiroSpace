package service

import (
	"easygin/internal/utils"
	"time"
)

func (s *Service) CreateUser(openId string) string {
	key := utils.GenerateRandomString(32)

	s.Redis.Set(openId, key, time.Duration(time.Minute*60))
	return key
}
