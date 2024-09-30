package model

type Report struct {
	ID          uint64
	Description string
	CreatedAt   int64 `gorm:"autoCreateTime"`
}
