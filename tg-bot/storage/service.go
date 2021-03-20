package storage

import (
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
)

type Repository interface {
	GetChatState(chatID int64) *entity.ChatState
	SetChatState(chatID int64, update *entity.ChatStateUpdate) (*entity.ChatState, error)
	GetSubscription(subID string) (*entity.Subscription, error)
	CreateSubscription(chatID int64, name string) (*entity.Subscription, error)
	UpdateSubscription(subID string, subUpdate *entity.SubscriptionUpdate) (*entity.Subscription, error)
	DeleteSubscription(subID string) error
	GetSubscriptions(chatID int64) []entity.Subscription
	GetEventSubscribers(location entity.EventData) (chatIDs []int64, err error)
}

type Service interface {
	GetChatState(chatID int64) *entity.ChatState
	SetChatState(chatID int64, update *entity.ChatStateUpdate) (*entity.ChatState, error)
	CreateSubscription(chatID int64, name string) (*entity.Subscription, error)
	GetSubscription(subID string) (*entity.Subscription, error)
	UpdateSubscription(subID string, subUpdate *entity.SubscriptionUpdate) (*entity.Subscription, error)
	DeleteSubscription(subID string) error
	GetSubscriptions(chatID int64) []entity.Subscription
	GetEventSubscribers(location entity.EventData) (chatIDs []int64, err error)

	SetAwaitUserInput(chatID int64, awaitInput string) error
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

func (s *service) GetSubscription(subID string) (*entity.Subscription, error) {
	return s.r.GetSubscription(subID)
}

func (s *service) UpdateSubscription(
	subID string, update *entity.SubscriptionUpdate,
) (*entity.Subscription, error) {
	return s.r.UpdateSubscription(subID, update)
}

func (s *service) CreateSubscription(chatID int64, name string) (*entity.Subscription, error) {
	return s.r.CreateSubscription(chatID, name)
}

func (s *service) DeleteSubscription(subID string) error {
	return s.r.DeleteSubscription(subID)
}

func (s *service) GetSubscriptions(chatID int64) []entity.Subscription {
	return s.r.GetSubscriptions(chatID)
}

func (s *service) SetAwaitUserInput(chatID int64, awaitInput string) error {
	stateUpdate := entity.ChatStateUpdate{AwaitInput: awaitInput}
	_, err := s.SetChatState(chatID, &stateUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetEventSubscribers(eventData entity.EventData) (chatIDs []int64, err error) {
	return s.r.GetEventSubscribers(eventData)
}
