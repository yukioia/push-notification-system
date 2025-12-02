package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// SubscriptionDB — структура для работы с таблицей subscriptions
type SubscriptionDB struct {
	ID       int64    `db:"id"`
	ClientID string   `db:"client_id"`
	Topics   []string `db:"topics"`
	Tags     []string `db:"tags"`
	Created  string   `db:"created_at"`
}

// PostgresSubscriptionRepo — реализация работы с таблицей subscriptions
type PostgresSubscriptionRepo struct {
	db *sqlx.DB
}

func NewPostgresSubscriptionRepo(db *sqlx.DB) *PostgresSubscriptionRepo {
	return &PostgresSubscriptionRepo{db: db}
}

func (r *PostgresSubscriptionRepo) CreateOrUpdate(s *SubscriptionDB) error {
	query := `
	INSERT INTO subscriptions (client_id, topics, tags)
	VALUES ($1, $2, $3)
	ON CONFLICT (client_id) DO UPDATE SET topics = EXCLUDED.topics, tags = EXCLUDED.tags
	RETURNING id, created_at;
	`
	return r.db.QueryRowx(query, s.ClientID, pq.Array(s.Topics), pq.Array(s.Tags)).Scan(&s.ID, &s.Created)
}

func (r *PostgresSubscriptionRepo) Delete(clientID string) error {
	_, err := r.db.Exec(`DELETE FROM subscriptions WHERE client_id=$1`, clientID)
	return err
}

func (r *PostgresSubscriptionRepo) GetByTopic(topic string) ([]SubscriptionDB, error) {
	var subs []SubscriptionDB
	err := r.db.Select(&subs, `SELECT id, client_id, topics, tags, created_at FROM subscriptions WHERE $1 = ANY(topics)`, topic)
	return subs, err
}
