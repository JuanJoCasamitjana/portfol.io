package model

type Report struct {
	ID          uint64
	Description string
	CreatedAt   RFC3339NanoTime `gorm:"autoCreateTime"`
}
