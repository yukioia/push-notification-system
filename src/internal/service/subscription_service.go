package service

import "push-server/src/internal/model"

// SubscriptionRepo — интерфейс для работы с подписками
type SubscriptionRepo interface {
	CreateOrUpdate(s *model.Subscription) error
	Delete(clientID string) error
	GetByTopic(topic string) ([]*model.Subscription, error)
	GetByClientID(clientID string) ([]*model.Subscription, error) // <--- добавлено
	GetClientIDsByTopic(topic string) ([]string, error)           // <--- добавлено
}

// SubscriptionService — бизнес-логика для подписок
type SubscriptionService struct {
	repo SubscriptionRepo
}

func NewSubscriptionService(r SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{repo: r}
}

func (s *SubscriptionService) Subscribe(sub *model.Subscription) error {
	return s.repo.CreateOrUpdate(sub)
}

func (s *SubscriptionService) Unsubscribe(clientID string) error {
	return s.repo.Delete(clientID)
}

// Получение подписок клиента
func (s *SubscriptionService) GetByClientID(clientID string) ([]*model.Subscription, error) {
	return s.repo.GetByClientID(clientID)
}

// Получение всех клиентов, подписанных на тему
func (s *SubscriptionService) GetClientsByTopic(topic string) ([]string, error) {
	return s.repo.GetClientIDsByTopic(topic)
}
