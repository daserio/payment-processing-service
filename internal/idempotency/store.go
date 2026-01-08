package idempotency

import "time"

type Result struct {
	TransactionID string
	Err           error
	CreatedAt     time.Time
}

type Store interface {
	Get(key string) (*Result, bool)
	Set(key string, result Result, ttl time.Duration)
}
