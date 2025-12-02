package service

// Service — главный контейнер сервисов
type Service struct {
	Subscription *SubscriptionService
	Message      *MessageService
}

// NewService — конструктор
func NewService(subRepo SubscriptionRepo, msgRepo MessageRepo) *Service {
	return &Service{
		Subscription: NewSubscriptionService(subRepo),
		Message:      NewMessageService(msgRepo),
	}
}
