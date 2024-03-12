package database

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"gorm.io/gorm"
)

func FindPostsPaginated(page, page_size int) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * page_size
	var posts []model.Post
	err := DB.Where("published = ?", true).Offset(offset).Limit(page_size).Find(&posts).Error
	return posts, err
}

func FindArticleByID(id uint64) (model.Article, error) {
	var article model.Article
	err := DB.First(&article, id).Error
	return article, err
}

func FindProjectByID(id uint64) (model.Project, error) {
	var project model.Project
	err := DB.First(&project, id).Error
	return project, err
}

func FindGalleryByID(id uint64) (model.Gallery, error) {
	var gallery model.Gallery
	err := DB.Preload("Images").First(&gallery, id).Error
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
	err := DB.Where("author = ?", author).Offset(offset).Limit(size).Find(&articles).Error
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
	err := DB.Model(&model.Gallery{}).Where("author = ?", author).Offset(offset).Limit(size).Preload("Images").
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

func FindTagLikeName(name string) ([]model.Tag, error) {
	var tags []model.Tag
	err := DB.Where("name LIKE ?", "%"+name+"%").Find(&tags).Error
	return tags, err
}

func DeleteGallery(gallery *model.Gallery) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(gallery).Delete(gallery).Error
	})
}
