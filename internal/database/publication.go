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

// Returns the collection with the given id and all its articles
func GetArticleCollectionById(id uint64) (model.ArticleCollection, error) {
	var collection model.ArticleCollection
	result := GormConn.Preload("Articles").First(&collection, id)
	return collection, result.Error
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

func GetAllImagesOfUser(username string) ([]model.Image, error) {
	var images []model.Image
	result := GormConn.Where("author = ?", username).Find(&images).Order("updated_at desc")
	return images, result.Error
}

func GetAllArticleCollectionsOfUser(username string) ([]model.ArticleCollection, error) {
	var collections []model.ArticleCollection
	result := GormConn.Where("author = ?", username).Find(&collections).Order("updated_at desc")
	return collections, result.Error
}

func GetAllImageCollectionsOfUser(username string) ([]model.ImageCollection, error) {
	var sesions []model.ImageCollection
	result := GormConn.Where("author = ?", username).Find(&sesions).Order("updated_at desc")
	return sesions, result.Error
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

func GetPaginatedImagesOfUser(username string, page int, pageSize int) ([]model.Image, error) {
	var images []model.Image
	result := GormConn.Where("author = ?", username).Order("updated_at desc").Limit(pageSize).Offset(page * pageSize).Find(&images)
	return images, result.Error
}

func GetPaginatedArticleCollectionsOfUser(username string, page int, pageSize int) ([]model.ArticleCollection, error) {
	var collections []model.ArticleCollection
	result := GormConn.Where("author = ?", username).Order("updated_at desc").Limit(pageSize).Offset(page * pageSize).Find(&collections)
	return collections, result.Error
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

// A pointer so collection can be modified and id can be used after creation
// Transaction so that if any of the operations fail, the whole transaction is rolled back
func CreateArticleCollection(collection *model.ArticleCollection) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(collection).Create(collection); result.Error != nil {
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

func AddArticleToArticleCollection(collection *model.ArticleCollection, article *model.Article) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(collection).Association("Articles").Append(article); result != nil {
			return result
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

func UpdateImage(image *model.Image) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(image).Updates(image); result.Error != nil {
			return result.Error
		}
		return nil
	})
}

func UpdateArticleCollection(collection *model.ArticleCollection) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(collection).Updates(collection); result.Error != nil {
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

func DeleteArticleFromArticleCollection(collection *model.ArticleCollection, article *model.Article) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(collection).Association("Articles").Delete(article); result != nil {
			return result
		}
		return nil
	})
}

func DeleteImageFromImageCollection(sesion *model.ImageCollection, image *model.Image) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Model(sesion).Association("Images").Delete(image); result != nil {
			return result
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

func GetPaginatedImages(page int, pageSize int) ([]model.Image, error) {
	var images []model.Image
	page = page - 1
	if page < 0 {
		page = 0
	}
	result := GormConn.Order("updated_at desc").Limit(pageSize).Offset(page * pageSize).Find(&images)
	return images, result.Error
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

func DeleteImageByID(id uint64) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Delete(&model.Image{}, id); result.Error != nil {
			return result.Error
		}
		return nil
	})
}
