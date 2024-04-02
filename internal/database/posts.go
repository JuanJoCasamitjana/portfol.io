package database

import (
	"errors"

	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"gorm.io/gorm"
)

func FindPostsPaginated(page, page_size int) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * page_size
	var posts []model.Post
	err := DB.Where("published = ?", true).Order("updated_at desc").Offset(offset).Limit(page_size).Find(&posts).Error
	return posts, err
}

func FindArticleByID(id uint64) (model.Article, error) {
	var article model.Article
	err := DB.Preload("Tags").First(&article, id).Error
	return article, err
}

func FindProjectByID(id uint64) (model.Project, error) {
	var project model.Project
	err := DB.Preload("Tags").First(&project, id).Error
	return project, err
}

func FindGalleryByID(id uint64) (model.Gallery, error) {
	var gallery model.Gallery
	err := DB.Preload("Images").Preload("Tags").First(&gallery, id).Error
	return gallery, err
}

func CreateArticle(article *model.Article) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Article{}).Create(article).Error
	})
}

func FindAllArticlesByAuthorPaginated(author string, page, size int) ([]model.Article, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	var articles []model.Article
	err := DB.Where("author = ?", author).Order("updated_at desc").Offset(offset).Limit(size).Find(&articles).Error
	return articles, err
}

func UpdateArticle(article *model.Article) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(article).Updates(article).Error
	})
}

func DeleteArticle(article *model.Article) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(article).Delete(article).Error
	})
}

func CreateGallery(gallery *model.Gallery) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Gallery{}).Create(gallery).Error
	})
}

func CreateImage(image *model.Image) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Image{}).Create(image).Error
	})
}

func UpdateGallery(gallery *model.Gallery) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(gallery).Updates(gallery).Error
	})
}

func DeleteImage(image *model.Image) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(image).Delete(image).Error
	})
}

func FindImageByID(id uint64) (model.Image, error) {
	var image model.Image
	err := DB.First(&image, id).Error
	return image, err
}

func FindAllGalleriesByAuthorPaginated(author string, page, size int) ([]model.Gallery, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	var galleries []model.Gallery
	err := DB.Model(&model.Gallery{}).Where("author = ?", author).Order("updated_at desc").Offset(offset).Limit(size).Preload("Images").
		Find(&galleries).Error
	return galleries, err
}

func AddTagToArticle(article *model.Article, tag *model.Tag) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(article).Association("Tags").Append(tag)
	})
}

func FindTagByName(name string) (model.Tag, error) {
	var tag model.Tag
	err := DB.Where("name = ?", name).First(&tag).Error
	return tag, err
}

func CreateTag(tag *model.Tag) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Tag{}).Create(tag).Error
	})
}

func RemoveTagFromArticle(article *model.Article, tag *model.Tag) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(article).Association("Tags").Delete(tag)
	})
}

func FindTagLikeName(name string, limit int) ([]model.Tag, error) {
	var tags []model.Tag
	err := DB.Where("name LIKE ?", "%"+name+"%").Limit(limit).Find(&tags).Error
	return tags, err
}

func DeleteGallery(gallery *model.Gallery) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(gallery).Delete(gallery).Error
	})
}

func FindPostableByTypeAndID(postableType string, id uint64) (model.Postable, error) {
	switch postableType {
	case "article":
		return FindArticleByID(id)
	case "project":
		return FindProjectByID(id)
	case "gallery":
		return FindGalleryByID(id)
	default:
		return nil, errors.New("invalid postable type")
	}
}

func FindAllArticlesByTagPaginated(tag string, page, size int) ([]model.Article, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	var articles []model.Article
	err := DB.Model(&model.Article{}).Joins("JOIN article_tags ON articles.id = article_tags.article_id").
		Joins("JOIN tags ON article_tags.tag_id = tags.id").Where("tags.name = ?", tag).Order("updated_at desc").
		Offset(offset).Limit(size).Find(&articles).Error
	return articles, err
}

func RemoveTagsFromArticle(tags []model.Tag, article *model.Article) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(article).Association("Tags").Delete(tags)
	})
}

func AddTagToGallery(gallery *model.Gallery, tag *model.Tag) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(gallery).Association("Tags").Append(tag)
	})
}

func RemoveTagsFromGallery(tags []model.Tag, gallery *model.Gallery) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(gallery).Association("Tags").Delete(tags)
	})
}

func FindAllGalleriesByTagPaginated(tag string, page, size int) ([]model.Gallery, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	var galleries []model.Gallery
	err := DB.Model(&model.Gallery{}).Preload("Images").Joins("JOIN gallery_tags ON galleries.id = gallery_tags.gallery_id").
		Joins("JOIN tags ON gallery_tags.tag_id = tags.id").Where("tags.name = ?", tag).Order("updated_at desc").
		Offset(offset).Limit(size).Find(&galleries).Error
	return galleries, err
}

func FindPostById(id uint64) (model.Post, error) {
	var post model.Post
	err := DB.First(&post, id).Error
	return post, err
}

func CountGalleries() (int64, error) {
	var count int64
	err := DB.Model(&model.Gallery{}).Count(&count).Error
	return count, err
}

func CountArticles() (int64, error) {
	var count int64
	err := DB.Model(&model.Article{}).Count(&count).Error
	return count, err
}

func FindPostsByQueryPaginated(query string, page, size int) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	var posts []model.Post
	err := DB.Where("title LIKE ? AND published = true", "%"+query+"%").Order("updated_at desc").Offset(offset).
		Limit(size).Find(&posts).Error
	return posts, err
}

func FindArticlesByQueryPaginated(query string, page, size int) ([]model.Article, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	var articles []model.Article
	err := DB.Where("title LIKE ? AND published = true", "%"+query+"%").Order("updated_at desc").Offset(offset).
		Limit(size).Find(&articles).Error
	return articles, err
}

func FindGalleriesByQueryPaginated(query string, page, size int) ([]model.Gallery, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	var galleries []model.Gallery
	err := DB.Where("title LIKE ?  AND published = true", "%"+query+"%").Order("updated_at desc").Preload("Images").
		Offset(offset).Limit(size).Find(&galleries).Error
	return galleries, err
}
