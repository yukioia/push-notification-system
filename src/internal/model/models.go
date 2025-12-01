package model

import "time"

type Subscription struct {
	ID       int64     `db:"id" json:"id"`
	ClientID string    `db:"client_id" json:"client_id"`
	Topics   []string  `db:"topics" json:"topics"`
	Tags     []string  `db:"tags" json:"tags"`
	Created  time.Time `db:"created_at" json:"created_at"`
}

type Message struct {
	ID      int64     `db:"id" json:"id"`
	Title   string    `db:"title" json:"title"`
	Body    string    `db:"body" json:"body"`
	Topic   string    `db:"topic" json:"topic"`
	Tags    []string  `db:"tags" json:"tags"`
	Created time.Time `db:"created_at" json:"created_at"`
}
