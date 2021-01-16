package storage

import "github.com/mightymatth/earthquake-tools/tg-bot/entity"

type Repository interface {
	GetChatState(chatID int64) *entity.ChatState
	SetChatState(chatID int64, update *entity.ChatStateUpdate) (*entity.ChatState, error)
	GetSubscription(chatID int64) *entity.Subscription
	SetSubscription(chatID int64, update *entity.SubscriptionUpdate) (*entity.Subscription, error)
	GetSubscriptions(chatID int64) []entity.Subscription
}

type Service interface {
	GetChatState(chatID int64) *entity.ChatState
	SetChatState(chatID int64, update *entity.ChatStateUpdate) (*entity.ChatState, error)
	GetSubscription(chatID int64) *entity.Subscription
	SetSubscription(chatID int64, update *entity.SubscriptionUpdate) (*entity.Subscription, error)
	GetSubscriptions(chatID int64) []entity.Subscription
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetChatState(chatID int64) *entity.ChatState {
	return s.r.GetChatState(chatID)
}

func (s *service) SetChatState(
	chatID int64, update *entity.ChatStateUpdate,
) (*entity.ChatState, error) {
	return s.r.SetChatState(chatID, update)
}

func (s *service) GetSubscription(chatID int64) *entity.Subscription {
	return s.r.GetSubscription(chatID)
}

func (s *service) SetSubscription(
	chatID int64, update *entity.SubscriptionUpdate,
) (*entity.Subscription, error) {
	return s.r.SetSubscription(chatID, update)
}

func (s *service) GetSubscriptions(chatID int64) []entity.Subscription {
	return s.r.GetSubscriptions(chatID)
}
