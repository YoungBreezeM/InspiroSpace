package service

import (
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

func (s *Service) SetToken(token string) (err error) {

	if _, err = s.Redis.RPush(TOKENS, token).Result(); err != nil {
		return
	}
	return
}

func (s *Service) GetToken() (token string, err error) {
	if token, err = s.Redis.LPop(TOKENS).Result(); err != nil {
		return
	}
	return
}

func (s *Service) SubChat(chat models.ApiChatRequest) (err error) {
	if _, err = s.Redis.Set(fmt.Sprintf("%s:%s", CHATQUEUE, chat.ChatId), chat, 0).Result(); err != nil {
		return
	}

	go func() {
		timer := time.NewTimer(1 * time.Second)
		<-timer.C

		if _, err = s.Redis.Del(fmt.Sprintf("%s:%s", CHATQUEUE, chat.ChatId)).Result(); err != nil {
			log.Error(err)
		}

		if err := s.SetToken(chat.Token); err != nil {
			log.Error(err)
		}
	}()

	return
}

func (s *Service) NotifyChat(chatId string) (chat *models.ApiChatRequest, err error) {
	chat = &models.ApiChatRequest{}
	if err = s.Redis.Get(fmt.Sprintf("%s:%s", CHATQUEUE, chatId)).Scan(&chat); err != nil {
		log.Error(err)
		return
	}
	return
}
