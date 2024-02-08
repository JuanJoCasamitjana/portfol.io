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
}

type Article struct {
	BasePost
	Title string `gorm:"not null"`
	Text  string `gorm:"not null"`
	Post  Post   `gorm:"polymorphic:Owner;"`
}

type Image struct {
	BasePost
	Title    string `gorm:"not null"`
	FilePath string `gorm:"not null"`
	Footer   string
	Post     Post `gorm:"polymorphic:Owner;"`
}

type ArticleCollection struct {
	BasePost
	Title    string
	Articles []Article `gorm:"many2many:collection_articles"`
}

type ImageCollection struct {
	BasePost
	Title  string
	Images []Image `gorm:"many2many:sesion_images"`
}

type Post struct {
	BasePost
	OwnerID   uint64
	OwnerType string
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

func (a *Article) AfterSave(tx *gorm.DB) error {
	var p Post
	p.OwnerType = "articles"
	p.OwnerID = a.ID
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&p).Error
	})
}

func (i *Image) AfterSave(tx *gorm.DB) error {
	var p Post
	p.OwnerType = "images"
	p.OwnerID = i.ID
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&p).Error
	})
}

func (a *Article) AfterDelete(tx *gorm.DB) error {
	var p Post
	p.OwnerType = "images"
	p.OwnerID = a.ID
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&p).Error
	})
}

func (i *Image) AfterDelete(tx *gorm.DB) error {
	var p Post
	p.OwnerType = "images"
	p.OwnerID = i.ID
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&p).Error
	})
}

func (a *Article) AfterUpdate(tx *gorm.DB) error {
	var p Post
	p.OwnerType = "images"
	p.OwnerID = a.ID
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&p).Error
	})
}

func (i *Image) AfterUpdate(tx *gorm.DB) error {
	var p Post
	p.OwnerType = "images"
	p.OwnerID = i.ID
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&p).Error
	})
}
