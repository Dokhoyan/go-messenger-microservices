package model

import "time"

type Record struct {
	ID        int64
	ChatID    int64
	Action    string
	CreatedAt time.Time
}