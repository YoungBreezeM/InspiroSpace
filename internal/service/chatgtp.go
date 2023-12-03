package service

import (
	"context"
	"easygin/internal/dao"
	"easygin/internal/log"
	"easygin/internal/models"
	"fmt"
	"time"
)

const TOKENS = "tokens"
const CHATQUEUE = "chat:queue"

func (s *Service) StoreChat(openId string, content string) (err error) {

	if err = s.Dao.CreateChat(&dao.ChatContent{OpenId: openId, Content: content}); err != nil {
		return
	}
	return
}

func (s *Service) SetToken(ctx context.Context, token string) (err error) {

	if _, err = s.Redis.RPush(ctx, TOKENS, token).Result(); err != nil {
		return
	}
	return
}

func (s *Service) GetToken(ctx context.Context) (token string, err error) {
	if token, err = s.Redis.LPop(ctx, TOKENS).Result(); err != nil {
		return
	}
	return
}

func (s *Service) SubChat(ctx context.Context, chat *models.ApiChatRequest) (err error) {
	go func() {
		timer := time.NewTimer(10 * time.Second)
		<-timer.C

		if _, err = s.Redis.Del(ctx, fmt.Sprintf("%s:%s", CHATQUEUE, chat.ChatId)).Result(); err != nil {
			log.Error(err)
		}

		if err := s.SetToken(ctx, chat.Token); err != nil {
			log.Error(err)
		}
	}()

	if _, err = s.Redis.Set(ctx, fmt.Sprintf("%s:%s", CHATQUEUE, chat.ChatId), chat, 0).Result(); err != nil {
		return
	}

	return
}

func (s *Service) NotifyChat(ctx context.Context, chatId string) (chat *models.ApiChatRequest, err error) {
	chat = &models.ApiChatRequest{}
	if err = s.Redis.Get(ctx, fmt.Sprintf("%s:%s", CHATQUEUE, chatId)).Scan(chat); err != nil {
		log.Error(err)
		return
	}
	return
}
