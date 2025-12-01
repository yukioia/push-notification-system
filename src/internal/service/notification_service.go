package service

import (
	"encoding/json"
	"push-server/src/internal/model"
	"push-server/src/internal/repo"
	"push-server/src/internal/sse"
)

type NotificationService struct {
	repo *repo.PostgresRepo
	sse  *sse.Manager
}

func NewNotificationService(r *repo.PostgresRepo, s *sse.Manager) *NotificationService {
	return &NotificationService{repo: r, sse: s}
}

func (ns *NotificationService) Subscribe(sub *model.Subscription) error {
	return ns.repo.CreateOrUpdateSubscription(sub)
}

func (ns *NotificationService) Unsubscribe(clientID string) error {
	// remove from db and close sse
	ns.sse.Remove(clientID)
	return ns.repo.DeleteSubscriptionByClientID(clientID)
}

func (ns *NotificationService) Publish(msg *model.Message) error {
	// save
	if err := ns.repo.SaveMessage(msg); err != nil {
		return err
	}
	// find matching subs
	subs, err := ns.repo.GetMatchingSubscriptions(msg.Topic, msg.Tags)
	if err != nil {
		return err
	}
	// prepare payload
	b, _ := json.Marshal(msg)
	// deliver to connected clients only (others will see message from DB on reconnect or via polling)
	var clientIDs []string
	for _, s := range subs {
		clientIDs = append(clientIDs, s.ClientID)
	}
	ns.sse.BroadcastToClients(clientIDs, string(b))
	return nil
}
