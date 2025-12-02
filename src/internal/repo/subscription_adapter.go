package repo

import (
	"push-server/src/internal/model"
	"time"
)

// Конвертеры между доменной моделью и DB
func toDB(s *model.Subscription) *SubscriptionDB {
	return &SubscriptionDB{
		ID:       s.ID,
		ClientID: s.ClientID,
		Topics:   s.Topics,
		Tags:     s.Tags,
		Created:  s.Created.Format(time.RFC3339),
	}
}

func fromDB(s *SubscriptionDB) *model.Subscription {
	t, _ := time.Parse(time.RFC3339, s.Created)
	return &model.Subscription{
		ID:       s.ID,
		ClientID: s.ClientID,
		Topics:   s.Topics,
		Tags:     s.Tags,
		Created:  t,
	}
}

// SubscriptionRepoAdapter реализует интерфейс service.SubscriptionRepo
type SubscriptionRepoAdapter struct {
	repo *PostgresSubscriptionRepo
}

func NewSubscriptionRepoAdapter(r *PostgresSubscriptionRepo) *SubscriptionRepoAdapter {
	return &SubscriptionRepoAdapter{repo: r}
}

func (a *SubscriptionRepoAdapter) CreateOrUpdate(s *model.Subscription) error {
	dbModel := toDB(s)
	if err := a.repo.CreateOrUpdate(dbModel); err != nil {
		return err
	}
	*s = *fromDB(dbModel)
	return nil
}

func (a *SubscriptionRepoAdapter) Delete(clientID string) error {
	return a.repo.Delete(clientID)
}

func (a *SubscriptionRepoAdapter) GetByTopic(topic string) ([]*model.Subscription, error) {
	dbSubs, err := a.repo.GetByTopic(topic)
	if err != nil {
		return nil, err
	}
	subs := make([]*model.Subscription, len(dbSubs))
	for i, s := range dbSubs {
		subs[i] = fromDB(&s)
	}
	return subs, nil
}

// ===== новые методы =====

// GetByClientID возвращает подписки конкретного клиента
func (a *SubscriptionRepoAdapter) GetByClientID(clientID string) ([]*model.Subscription, error) {
	dbSubs, err := a.repo.GetByClientID(clientID)
	if err != nil {
		return nil, err
	}
	subs := make([]*model.Subscription, len(dbSubs))
	for i, s := range dbSubs {
		subs[i] = fromDB(&s)
	}
	return subs, nil
}

// GetClientIDsByTopic возвращает список client_id, подписанных на тему
func (a *SubscriptionRepoAdapter) GetClientIDsByTopic(topic string) ([]string, error) {
	return a.repo.GetClientIDsByTopic(topic)
}
