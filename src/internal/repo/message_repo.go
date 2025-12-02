package repo

import (
	"github.com/lib/pq"
	"push-server/src/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

// MessageDB — структура для работы с таблицей messages
type MessageDB struct {
	ID      int64    `db:"id"`
	Title   string   `db:"title"`
	Body    string   `db:"body"`
	Topic   string   `db:"topic"`
	Tags    []string `db:"tags"`
	Created string   `db:"created_at"`
}

// PostgresMessageRepo — реализация работы с таблицей messages
type PostgresMessageRepo struct {
	db *sqlx.DB
}

func NewPostgresMessageRepo(db *sqlx.DB) *PostgresMessageRepo {
	return &PostgresMessageRepo{db: db}
}

func (r *PostgresMessageRepo) Save(m *MessageDB) error {
	query := `
	INSERT INTO messages (title, body, topic, tags)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at;
	`
	return r.db.QueryRowx(query, m.Title, m.Body, m.Topic, pq.Array(m.Tags)).Scan(&m.ID, &m.Created)
}

func (r *PostgresMessageRepo) GetByTopic(topic string) ([]MessageDB, error) {
	var msgs []MessageDB
	err := r.db.Select(&msgs, `SELECT id, title, body, topic, tags, created_at FROM messages WHERE topic=$1`, topic)
	return msgs, err
}

// ---- Конвертеры между доменной моделью и DB ----
func messageToDB(m *model.Message) *MessageDB {
	return &MessageDB{
		ID:      m.ID,
		Title:   m.Title,
		Body:    m.Body,
		Topic:   m.Topic,
		Tags:    m.Tags,
		Created: m.Created.Format(time.RFC3339),
	}
}

func messageFromDB(m *MessageDB) *model.Message {
	t, _ := time.Parse(time.RFC3339, m.Created)
	return &model.Message{
		ID:      m.ID,
		Title:   m.Title,
		Body:    m.Body,
		Topic:   m.Topic,
		Tags:    m.Tags,
		Created: t,
	}
}

// MessageRepoAdapter — адаптер для сервиса
type MessageRepoAdapter struct {
	repo *PostgresMessageRepo
}

func NewMessageRepoAdapter(r *PostgresMessageRepo) *MessageRepoAdapter {
	return &MessageRepoAdapter{repo: r}
}

func (a *MessageRepoAdapter) Save(m *model.Message) error {
	dbModel := messageToDB(m)
	if err := a.repo.Save(dbModel); err != nil {
		return err
	}
	*m = *messageFromDB(dbModel)
	return nil
}

func (a *MessageRepoAdapter) GetByTopic(topic string) ([]*model.Message, error) {
	dbMsgs, err := a.repo.GetByTopic(topic)
	if err != nil {
		return nil, err
	}
	msgs := make([]*model.Message, len(dbMsgs))
	for i, m := range dbMsgs {
		msgs[i] = messageFromDB(&m)
	}
	return msgs, nil
}
