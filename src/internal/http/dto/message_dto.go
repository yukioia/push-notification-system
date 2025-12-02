package dto

import "push-server/src/internal/model"

type MessageDTO struct {
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Topic string   `json:"topic"`
	Tags  []string `json:"tags"`
}

// MessageDTOFromDomain конвертирует модель в DTO
func MessageDTOFromDomain(m *model.Message) *MessageDTO {
	return &MessageDTO{
		Title: m.Title,
		Body:  m.Body,
		Topic: m.Topic,
		Tags:  m.Tags,
	}
}
