// internal/model/subscription.go
package model

import "time"

type Subscription struct {
	ID       int64
	ClientID string
	Topics   []string
	Tags     []string
	Created  time.Time
}
