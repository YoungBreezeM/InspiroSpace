package service

import (
	"easygin/internal/dao"
)

func (s *Service) StoreChat(openId string, content string) (err error) {

	if err = s.Dao.CreateChat(&dao.ChatContent{OpenId: openId, Content: content}); err != nil {
		return
	}
	return
}
