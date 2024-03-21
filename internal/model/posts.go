package model

import (
	"time"

	"gorm.io/gorm"
)

type BasePost struct {
	ID        uint64
	Title     string
	Author    string
	User      User `gorm:"foreignKey:Author;references:Username"`
	Published bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tag struct {
	ID   uint64
	Name string `gorm:"unique"`
}

type Article struct {
	BasePost
	Tags    []Tag `gorm:"many2many:article_tags;"`
	Content string
}

type Project struct {
	BasePost
	Description string
	Tags        []Tag `gorm:"many2many:project_tags;"`
	Link        string
}

type Image struct {
	ID        uint64
	Owner     string
	User      User `gorm:"foreignKey:Owner;references:Username"`
	Footer    string
	ImageURL  string
	ThumbURL  string
	DeleteURL string
	GalleryID uint64
	Gallery   Gallery
}

type Gallery struct {
	BasePost
	Tags   []Tag `gorm:"many2many:gallery_tags;"`
	Images []Image
}

type Post struct {
	BasePost
	Tags      []Tag `gorm:"many2many:post_tags;"`
	OwnerID   uint64
	OwnerType string
}

// Postable is an interface that all posts must implement
// It allows to easily index posts in a single slice
type Postable interface {
	GetID() uint64
	GetTitle() string
	GetAuthor() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetTags() []Tag
}

func (p BasePost) GetID() uint64 {
	return p.ID
}

func (p BasePost) GetTitle() string {
	return p.Title
}

func (p BasePost) GetAuthor() string {
	return p.Author
}

func (p BasePost) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p BasePost) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

func (a Article) GetTags() []Tag {
	return a.Tags
}

func (p Project) GetTags() []Tag {
	return p.Tags
}

func (g Gallery) GetTags() []Tag {
	return g.Tags
}

func (p Post) GetTags() []Tag {
	return p.Tags
}

// AfterCreate is a hook that creates a post after creating an article
// It helps indexing posts arbitrarily
func (a *Article) AfterCreate(tx *gorm.DB) error {
	var post Post
	post.OwnerID = a.ID
	post.OwnerType = "article"
	post.Author = a.Author
	post.Title = a.Title
	post.Published = a.Published
	tx.Transaction(func(tx *gorm.DB) error {
		tx.Create(&post)
		return nil
	})
	return nil
}

func (p *Project) AfterCreate(tx *gorm.DB) error {
	var post Post
	post.OwnerID = p.ID
	post.OwnerType = "project"
	post.Author = p.Author
	post.Title = p.Title
	post.Published = p.Published
	tx.Transaction(func(tx *gorm.DB) error {
		tx.Create(&post)
		return nil
	})
	return nil
}

func (g *Gallery) AfterCreate(tx *gorm.DB) error {
	var post Post
	post.OwnerID = g.ID
	post.OwnerType = "gallery"
	post.Author = g.Author
	post.Title = g.Title
	post.Published = g.Published
	tx.Transaction(func(tx *gorm.DB) error {
		tx.Create(&post)
		return nil
	})
	return nil
}

func (a *Article) AfterSave(tx *gorm.DB) error {
	var post Post
	tx.Where("owner_id = ? AND owner_type = ?", a.ID, "article").Preload("Tags").First(&post)
	post.Title = a.Title
	post.Author = a.Author
	post.Published = a.Published
	post.Tags = a.Tags
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&post).Error
	})
}

func (p *Project) AfterSave(tx *gorm.DB) error {
	var post Post
	tx.Where("owner_id = ? AND owner_type = ?", p.ID, "project").Preload("Tags").First(&post)
	post.Title = p.Title
	post.Author = p.Author
	post.Published = p.Published
	post.Tags = p.Tags
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&post).Error
	})
}

func (g *Gallery) AfterSave(tx *gorm.DB) error {
	var post Post
	tx.Where("owner_id = ? AND owner_type = ?", g.ID, "gallery").Preload("Tags").First(&post)
	post.Title = g.Title
	post.Author = g.Author
	post.Published = g.Published
	post.Tags = g.Tags
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&post).Error
	})
}

func (a *Article) BeforeDelete(tx *gorm.DB) error {
	var post Post
	tx.Where("owner_id = ? AND owner_type = ?", a.ID, "article").First(&post)
	return tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&post).Association("Tags").Clear()
		if err != nil {
			return err
		}
		return tx.Delete(&post).Error
	})
}

func (p *Project) BeforeDelete(tx *gorm.DB) error {
	var post Post
	tx.Where("owner_id = ? AND owner_type = ?", p.ID, "project").First(&post)
	return tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&post).Association("Tags").Clear()
		if err != nil {
			return err
		}
		return tx.Delete(&post).Error
	})
}

func (g *Gallery) BeforeDelete(tx *gorm.DB) error {
	var post Post
	tx.Where("owner_id = ? AND owner_type = ?", g.ID, "gallery").First(&post)
	return tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&post).Association("Tags").Clear()
		if err != nil {
			return err
		}
		return tx.Delete(&post).Error
	})
}

func (t *Tag) ColorOfTag() string {
	colors := []string{"#C84630", "#FFB627", "#219797", "#6113CD", "#1A5E63"}
	sum := 0
	for i := range t.Name {
		sum += int(t.Name[i])
	}
	return colors[sum%5]
}
