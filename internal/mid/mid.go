package mid

import "easygin/internal/service"

type Mid struct {
	S *service.Service
}

func NewMid(s *service.Service) *Mid {
	return &Mid{
		s,
	}
}
