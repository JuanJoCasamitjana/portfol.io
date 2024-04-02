package model

import "time"

type Report struct {
	ID          uint64
	Description string
	CreatedAt   time.Time
}
