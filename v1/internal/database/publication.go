package database

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"gorm.io/gorm"
)

func GetArticleById(id uint64) (model.Article, error) {
	var article model.Article
	result := GormConn.First(&article, id)
	return article, result.Error
}

func GetImageById(id uint64) (model.Image, error) {
	var image model.Image
	result := GormConn.First(&image, id)
	return image, result.Error
}

// Returns the sesion with the given id and all its images
func GetImageCollectionById(id uint64) (model.ImageCollection, error) {
	var sesion model.ImageCollection
	result := GormConn.Preload("Images").First(&sesion, id)
	return sesion, result.Error
}

func GetImageCollectionByIdIfPublished(id uint64) (model.ImageCollection, error) {
	var sesion model.ImageCollection
	result := GormConn.Preload("Images").Where("published = ?", true).First(&sesion, id)
	return sesion, result.Error
}

func GetAllArticlesOfUser(username string) ([]model.Article, error) {
	var articles []model.Article
	result := GormConn.Where("author = ?", username).Find(&articles).Order("updated_at desc")
	return articles, result.Error
}

func GetAllImageCollectionsOfUserWithImages(username string) ([]model.ImageCollection, error) {
	var sesions []model.ImageCollection
	result := GormConn.Preload("Images").Where("author = ?", username).Find(&sesions).Order("updated_at desc")
	return sesions, result.Error
}

func GetPaginatedArticlesOfUser(username string, page int, pageSize int) ([]model.Article, error) {
	var articles []model.Article
	result := GormConn.Where("author = ?", username).Order("updated_at desc").Limit(pageSize).Offset(page * pageSize).Find(&articles)
	return articles, result.Error
}

func GetPaginatedImageCollectionsOfUser(username string, page int, pageSize int) ([]model.ImageCollection, error) {
	var sesions []model.ImageCollection
	result := GormConn.Where("author = ?", username).Order("updated_at desc").Limit(pageSize).Offset(page * pageSize).Find(&sesions)
	return sesions, result.Error
}

// A pointer so article can be modified and id can be used after creation
// Transaction so that if any of the operations fail, the whole transaction is rolled back
func CreateArticle(article *model.Article) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(article).Create(article); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

// A pointer so image can be modified and id can be used after creation
// Transaction so that if any of the operations fail, the whole transaction is rolled back
func CreateImage(image *model.Image) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(image).Create(image); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

// A pointer so sesion can be modified and id can be used after creation
// Transaction so that if any of the operations fail, the whole transaction is rolled back
func CreateImageCollection(sesion *model.ImageCollection) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(sesion).Create(sesion); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func AddImageToImageCollection(sesion *model.ImageCollection, image *model.Image) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(sesion).Association("Images").Append(image); result != nil {
			return result
		}
		return nil
	})
}

func UpdateArticle(article *model.Article) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(article).Updates(article); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func UpdateImageCollection(sesion *model.ImageCollection) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(sesion).Updates(sesion); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func DeleteArticle(article *model.Article) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(article).Delete(article); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func DeleteImage(image *model.Image) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(image).Delete(image); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func DeleteArticleCollection(collection *model.ArticleCollection) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(collection).Delete(collection); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func DeleteImageCollection(sesion *model.ImageCollection) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(sesion).Delete(sesion); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func GetPaginatedArticles(page int, pageSize int) ([]model.Article, error) {
	var articles []model.Article
	page = page - 1
	if page < 0 {
		page = 0
	}
	result := GormConn.Order("updated_at desc").Limit(pageSize).Offset(page * pageSize).Find(&articles)
	return articles, result.Error
}

func GetPostsOfUser(username string, page int, pageSize int) ([]model.Post, error) {
	var posts []model.Post
	page = page - 1
	if page < 0 {
		page = 0
	}
	result := GormConn.Where("author = ?", username).Order("updated_at desc").Limit(pageSize).
		Offset(page * pageSize).Preload("User").Preload("BasePost").Find(&posts)
	return posts, result.Error
}

func GetPosts(page int, pageSize int) ([]model.Post, error) {
	var posts []model.Post
	page = page - 1
	if page < 0 {
		page = 0
	}
	result := GormConn.Order("updated_at desc").Limit(pageSize).Offset(page * pageSize).
		Preload("User").Preload("BasePost").Find(&posts)
	return posts, result.Error
}

func CreateTag(tag *model.Tag) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(tag).Create(tag); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func GetTagById(id uint64) (model.Tag, error) {
	var tag model.Tag
	result := GormConn.First(&tag, id)
	return tag, result.Error
}

func GetTagByName(name string) (model.Tag, error) {
	var tag model.Tag
	result := GormConn.Where("name = ?", name).First(&tag)
	return tag, result.Error
}

func GetAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	result := GormConn.Find(&tags)
	return tags, result.Error
}

func DeleteTag(tag *model.Tag) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(tag).Delete(tag); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func GetTagsOfArticle(article *model.Article) ([]model.Tag, error) {
	var tags []model.Tag
	err := GormConn.Model(article).Association("Tags").Find(&tags)
	return tags, err
}

func GetTagsOfImage(image *model.Image) ([]model.Tag, error) {
	var tags []model.Tag
	err := GormConn.Model(image).Association("Tags").Find(&tags)
	return tags, err
}
