// internal/model/message.go
package model

import "time"

type Message struct {
	ID      int64
	Title   string
	Body    string
	Topic   string
	Tags    []string
	Created time.Time
}
