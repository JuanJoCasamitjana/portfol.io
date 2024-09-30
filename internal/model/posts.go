package model

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/JuanJoCasamitjana/portfol.io/internal/utils"
	"gorm.io/gorm"
)

type BasePost struct {
	ID        uint64
	Title     string
	Author    string
	User      User `gorm:"foreignKey:Author;references:Username"`
	Published bool
	CreatedAt ISOTime `gorm:"autoCreateTime"`
	UpdatedAt ISOTime `gorm:"autoUpdateTime"`
}

type Tag struct {
	ID   uint64
	Name string `gorm:"unique"`
}

// A user can vote for a specific tag on a post
type Vote struct {
	ID    uint64
	Voter string
	User  User `gorm:"foreignKey:Voter;references:Username"`
	TagID uint64
	Tag   Tag `gorm:"foreignKey:TagID;references:ID"`
}

type Article struct {
	BasePost
	Votes   []Vote `gorm:"many2many:article_votes;"`
	Content string
}

type Project struct {
	BasePost
	Description string
	Votes       []Vote `gorm:"many2many:project_votes;"`
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
	Votes  []Vote `gorm:"many2many:gallery_votes;"`
	Images []Image
}

type Post struct {
	BasePost
	Votes     []Vote `gorm:"many2many:post_votes;"`
	OwnerID   uint64
	OwnerType string
}

// Postable is an interface that all posts must implement
// It allows to easily index posts in a single slice
type Postable interface {
	GetID() uint64
	GetTitle() string
	GetAuthor() string
	GetCreatedAt() ISOTime
	GetUpdatedAt() ISOTime
	GetVotes() []Vote
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

func (p BasePost) GetCreatedAt() ISOTime {
	return p.CreatedAt
}

func (p BasePost) GetUpdatedAt() ISOTime {
	return p.UpdatedAt
}

func (a Article) GetVotes() []Vote {
	return a.Votes
}

func (p Project) GetVotes() []Vote {
	return p.Votes
}

func (g Gallery) GetVotes() []Vote {
	return g.Votes
}

func (p Post) GetVotes() []Vote {
	return p.Votes
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
	tx.Where("owner_id = ? AND owner_type = ?", a.ID, "article").Preload("Votes").First(&post)
	post.Title = a.Title
	post.Author = a.Author
	post.Published = a.Published
	post.Votes = a.Votes
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&post).Error
	})
}

func (p *Project) AfterSave(tx *gorm.DB) error {
	var post Post
	tx.Where("owner_id = ? AND owner_type = ?", p.ID, "project").Preload("Votes").First(&post)
	post.Title = p.Title
	post.Author = p.Author
	post.Published = p.Published
	post.Votes = p.Votes
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&post).Error
	})
}

func (g *Gallery) AfterSave(tx *gorm.DB) error {
	var post Post
	tx.Where("owner_id = ? AND owner_type = ?", g.ID, "gallery").Preload("Votes").First(&post)
	post.Title = g.Title
	post.Author = g.Author
	post.Published = g.Published
	post.Votes = g.Votes
	return tx.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&post).Error
	})
}

func (p *Post) AfterCreate(tx *gorm.DB) error {
	owner := p.User
	var users_to_notify []User
	tx.Model(&User{}).Where("username IN (SELECT username FROM follows WHERE owner = ?)", owner.Username).
		Find(&users_to_notify)
	var emails []string
	for _, user := range users_to_notify {
		emails = append(emails, user.Email)
	}
	data := map[string]any{
		"title":  p.Title,
		"author": p.Author,
		"type":   p.OwnerType,
	}
	go sendNotification(emails, data)
	return nil
}

func (a *Article) BeforeDelete(tx *gorm.DB) error {
	var post Post
	tx.Where("owner_id = ? AND owner_type = ?", a.ID, "article").First(&post)
	return tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&post).Association("Votes").Clear()
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
		var sections []Section
		err := tx.Model(post).Association("Sections").Find(&sections)
		if err != nil {
			return err
		}
		for _, section := range sections {
			err = tx.Model(section).Association("Posts").Delete(post)
			if err != nil {
				return err
			}
		}
		err = tx.Model(&post).Association("Votes").Clear()
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
		err := tx.Model(&post).Association("Votes").Clear()
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

func sendNotification(emails []string, data map[string]any) {
	var body bytes.Buffer
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := fmt.Sprintf("Subject: New %s by %s\n%s\n\n", data["type"], data["author"], headers)
	body.WriteString(subject)
	t := template.Must(template.ParseFiles("web/templates/email_notification.html"))
	t.Execute(&body, data)
	utils.SendEmailNotification(emails, body.Bytes())
}
