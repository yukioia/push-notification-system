package service

import "push-server/src/internal/model"

// MessageRepo — интерфейс репозитория сообщений
type MessageRepo interface {
	Save(msg *model.Message) error
	GetByTopic(topic string) ([]*model.Message, error)
}

// MessageService — бизнес-логика для сообщений
type MessageService struct {
	repo MessageRepo
}

func NewMessageService(repo MessageRepo) *MessageService {
	return &MessageService{repo: repo}
}

func (ms *MessageService) Publish(msg *model.Message) error {
	return ms.repo.Save(msg)
}

func (ms *MessageService) GetMessagesByTopic(topic string) ([]*model.Message, error) {
	return ms.repo.GetByTopic(topic)
}
