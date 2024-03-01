package model

type Profile struct {
	ID        uint64
	UserID    uint64
	User      User
	Bio       string
	ImageURL  string
	ThumbURL  string
	DeleteURL string
	CreatedAt string
	UpdatedAt string
}

type Section struct {
	ID        uint64
	ProfileID uint64
	Profile   Profile
	Title     string
	Posts     []Post `gorm:"many2many:section_posts;"`
}
