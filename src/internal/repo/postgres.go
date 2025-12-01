package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"push-server/src/internal/model"
)

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(dsn string) (*PostgresRepo, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{db: db}, nil
}

func (r *PostgresRepo) CreateOrUpdateSubscription(s *model.Subscription) error {
	// upsert by client_id
	q := `
    INSERT INTO subscriptions (client_id, topics, tags)
    VALUES ($1, $2, $3)
    ON CONFLICT (client_id) DO UPDATE SET topics = EXCLUDED.topics, tags = EXCLUDED.tags
    RETURNING id, created_at;
    `
	return r.db.QueryRowx(q, s.ClientID, pqStringArray(s.Topics), pqStringArray(s.Tags)).Scan(&s.ID, &s.Created)
}

func (r *PostgresRepo) DeleteSubscriptionByClientID(clientID string) error {
	_, err := r.db.Exec(`DELETE FROM subscriptions WHERE client_id=$1`, clientID)
	return err
}

func (r *PostgresRepo) GetMatchingSubscriptions(topic string, tags []string) ([]model.Subscription, error) {
	// simple matching: topic contained in topics array OR topics empty -> subscribe to all
	// and if tags provided, check overlap
	var subs []model.Subscription
	query := `
    SELECT id, client_id, topics, tags, created_at FROM subscriptions
    WHERE ($1 = ANY(topics) OR array_length(topics,1) IS NULL)
    `
	err := r.db.Select(&subs, query, topic)
	if err != nil {
		return nil, err
	}
	// server-side tag overlap filtering (simple)
	if len(tags) == 0 {
		return subs, nil
	}
	var matched []model.Subscription
	tagSet := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		tagSet[t] = struct{}{}
	}
	for _, s := range subs {
		for _, st := range s.Tags {
			if _, ok := tagSet[st]; ok {
				matched = append(matched, s)
				break
			}
		}
	}
	return matched, nil
}

func (r *PostgresRepo) SaveMessage(m *model.Message) error {
	q := `INSERT INTO messages (title, body, topic, tags) VALUES ($1,$2,$3,$4) RETURNING id, created_at`
	return r.db.QueryRowx(q, m.Title, m.Body, m.Topic, pqStringArray(m.Tags)).Scan(&m.ID, &m.Created)
}

func pqStringArray(s []string) interface{} {
	if s == nil {
		return []string{}
	}
	return pq.Array(s)
}
