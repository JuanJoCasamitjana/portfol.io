package model

type Report struct {
	ID          uint64
	Description string
	CreatedAt   ISOTime `gorm:"autoCreateTime"`
}
