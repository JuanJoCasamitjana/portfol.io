package model

import (
	"time"

	"gorm.io/gorm"
)

// This acts like an abdstract class but not quite so
// It's a struct that contains the common fields of Article and Image
// It can still be instantiated, but it's not meant to be
type BasePost struct {
	ID        uint64    `gorm:"primaryKey"`
	Author    string    `gorm:"not null"`
	User      User      `gorm:"foreignKey:Author;references:Username"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Tags      []Tag     `gorm:"many2many:post_tags;"`
}

type Article struct {
	BasePost
	Title string `gorm:"not null"`
	Text  string `gorm:"not null"`
	Post  Post   `gorm:"polymorphic:Owner;"`
}

type Image struct {
	BasePost
	FilePath  string `gorm:"not null"`
	ThumbPath string `gorm:"not null"`
	DeleteUrl string `gorm:"not null"`
	Footer    string
	Post      Post `gorm:"polymorphic:Owner;"`
}

type ArticleCollection struct {
	BasePost
	Title    string
	Articles []Article `gorm:"many2many:collection_articles"`
}

type ImageCollection struct {
	BasePost
	Title     string
	Images    []Image `gorm:"many2many:sesion_images"`
	Published bool    `gorm:"default:false"`
}

type Post struct {
	BasePost
	OwnerID   uint64
	OwnerType string
}

type Tag struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string `gorm:"not null;unique"`
}

type Postable interface {
	GetBasePost() BasePost
}

func (a Article) GetBasePost() *BasePost {
	return &a.BasePost
}

func (i Image) GetBasePost() *BasePost {
	return &i.BasePost
}

func (ic ImageCollection) GetBasePost() *BasePost {
	return &ic.BasePost
}

func (a *Article) AfterCreate(tx *gorm.DB) error {
	var p Post
	p.OwnerType = "articles"
	p.OwnerID = a.ID
	p.Author = a.Author
	p.Tags = a.Tags
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&p).Error
	})
}

func (ic *ImageCollection) AfterCreate(tx *gorm.DB) error {
	var p Post
	p.OwnerType = "image_collections"
	p.OwnerID = ic.ID
	p.Author = ic.Author
	p.Tags = ic.Tags
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&p).Error
	})
}

func (a *Article) AfterSave(tx *gorm.DB) error {
	var p Post
	err := tx.Where("owner_id = ?", a.ID).First(&p).Error
	if err != nil {
		return err
	}
	p.OwnerType = "articles"
	p.OwnerID = a.ID
	p.Author = a.Author
	p.Tags = a.Tags
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&p).Error
	})
}

func (ic *ImageCollection) AfterSave(tx *gorm.DB) error {
	var p Post
	err := tx.Where("owner_id = ?", ic.ID).First(&p).Error
	if err != nil {
		return err
	}
	p.OwnerType = "image_collections"
	p.OwnerID = ic.ID
	p.Author = ic.Author
	p.Tags = ic.Tags
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&p).Error
	})
}

func (a *Article) AfterDelete(tx *gorm.DB) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Where("owner_id = ?", a.ID).Delete(&Post{}).Error
	})
}

func (ic *ImageCollection) AfterDelete(tx *gorm.DB) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Where("owner_id = ?", ic.ID).Delete(&Post{}).Error
	})
}
