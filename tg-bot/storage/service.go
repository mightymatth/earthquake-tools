package storage

import "github.com/mightymatth/earthquake-tools/tg-bot/entity"

type Repository interface {
	GetChatState(chatID string) (*entity.ChatState, error)
}

type Service interface {
	GetChatState(chatID string) (*entity.ChatState, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetChatState(chatID string) (*entity.ChatState, error)  {
	return s.r.GetChatState(chatID)
}
