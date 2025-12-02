// internal/http/dto/subscription_dto.go
package dto

import "push-server/src/internal/model"

type SubscriptionDTO struct {
	ClientID string   `json:"client_id"`
	Topics   []string `json:"topics"`
	Tags     []string `json:"tags"`
}

func FromDomain(s *model.Subscription) *SubscriptionDTO {
	return &SubscriptionDTO{
		ClientID: s.ClientID,
		Topics:   s.Topics,
		Tags:     s.Tags,
	}
}
