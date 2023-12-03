package service

import (
	"context"
	"easygin/internal/utils"
	"time"
)

func (s *Service) CreateUser(ctx context.Context, openId string) string {
	key := utils.GenerateRandomString(32)

	s.Redis.Set(ctx, openId, key, time.Duration(time.Minute*60))
	return key
}
