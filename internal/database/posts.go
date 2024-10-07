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
	err := DB.Preload("Votes.Tag").First(&article, id).Error
	return article, err
}

func FindProjectByID(id uint64) (model.Project, error) {
	var project model.Project
	err := DB.Preload("Votes.Tag").First(&project, id).Error
	return project, err
}

func FindGalleryByID(id uint64) (model.Gallery, error) {
	var gallery model.Gallery
	err := DB.Preload("Images").Preload("Votes.Tag").First(&gallery, id).Error
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

func FindTagLikeName(name string, limit int) ([]model.Tag, error) {
	var tags []model.Tag
	err := DB.Where("name LIKE ?", "%"+name+"%").Limit(limit).Find(&tags).Error
	return tags, err
}

func DeleteGallery(gallery *model.Gallery) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(gallery).Association("Images").Clear()
		if err != nil {
			return err
		}
		return tx.Model(gallery).Delete(gallery).Error
	})
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

func FindAllPostsPaginated(page, size int) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	var posts []model.Post
	err := DB.Order("updated_at desc").Offset(offset).Limit(size).Find(&posts).Error
	return posts, err
}

func FindAllPostsByqueryPaginated(page, size int, query string) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	var posts []model.Post
	err := DB.Where("title LIKE ?", "%"+query+"%").Order("updated_at desc").Offset(offset).Limit(size).Find(&posts).Error
	return posts, err
}

func DeletePostByID(id uint64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var post model.Post
		err := tx.First(&post, id).Error
		if err != nil {
			return err
		}
		if post.OwnerType == "article" {
			var article model.Article
			err = tx.First(&article, post.OwnerID).Error
			if err != nil {
				return err
			}
			return tx.Model(&article).Delete(&article).Error
		}
		if post.OwnerType == "gallery" {
			var gallery model.Gallery
			err = tx.First(&gallery, post.OwnerID).Error
			if err != nil {
				return err
			}
			return tx.Model(&gallery).Delete(&gallery).Error
		}
		return errors.New("invalid post owner type")
	})
}

func FilterPostsInUserSection(posts []model.Post, username, section string) ([]model.Post, error) {
	var allIDs []uint64
	var filteredIDs []uint64
	var sectionDB model.Section
	var filteredPosts []model.Post
	for _, post := range posts {
		allIDs = append(allIDs, post.ID)
	}
	err := DB.Where("owner = ? AND name = ?", username, section).First(&sectionDB).Error
	if err != nil {
		return nil, err
	}
	err = DB.Table("section_posts").Where("section_id = ? AND post_id IN (?)", sectionDB.ID, allIDs).Pluck("post_id", &filteredIDs).Error
	if err != nil {
		return nil, err
	}
	for _, post := range posts {
		for _, id := range filteredIDs {
			if post.ID == id {
				filteredPosts = append(filteredPosts, post)
			}
		}
	}
	return filteredPosts, nil
}

func GetFirstFiftyMostVotedTagsForArticle(articleID uint64) ([]model.Tag, error) {
	var tags []model.Tag
	err := DB.Table("tags").
		Select("tags.id, tags.name").Joins("JOIN votes ON votes.tag_id = tags.id").
		Joins("JOIN article_votes ON article_votes.vote_id = votes.id").
		Where("article_votes.article_id = ?", articleID).Group("tags.id, tags.name").
		Order("COUNT(votes.id) DESC").Limit(50).Scan(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func GetFirstFiftyMostVotedTagsForGallery(galleryID uint64) ([]model.Tag, error) {
	var tags []model.Tag
	err := DB.Table("tags").
		Select("tags.id, tags.name").Joins("JOIN votes ON votes.tag_id = tags.id").
		Joins("JOIN gallery_votes ON gallery_votes.vote_id = votes.id").
		Where("gallery_votes.article_id = ?", galleryID).Group("tags.id, tags.name").
		Order("COUNT(votes.id) DESC").Limit(50).Scan(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func VoteTagForArticle(article *model.Article, vote *model.Vote) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(article).Association("Votes").Append(vote)
	})
}

func VoteTagForGallery(gallery *model.Gallery, vote *model.Vote) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(gallery).Association("Votes").Append(vote)
	})
}

func UnvoteTagForArticle(article *model.Article, vote *model.Vote) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(article).Association("Votes").Delete(vote)
	})
}

func UnvoteTagForGallery(gallery *model.Gallery, vote *model.Vote) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(gallery).Association("Votes").Delete(vote)
	})
}
func RemoveAllVotesForArticle(article *model.Article) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(article).Association("Votes").Clear()
	})
}

func RemoveAllVotesForGallery(gallery *model.Gallery) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(gallery).Association("Votes").Clear()
	})
}

func FindVoteByTagAndUser(tagID uint64, username string) (model.Vote, error) {
	var vote model.Vote
	err := DB.Model(vote).Where("tag_id = ? AND voter = ?", tagID, username).First(&vote).Error
	return vote, err
}

func VoteExistsForTagUserAndPost(tagID uint64, voter string, postID uint64, postType string) bool {
	var count int64
	err := DB.Table("votes").Joins("JOIN "+postType+"_votes ON votes.id = "+postType+"_votes.vote_id").
		Where("votes.tag_id = ? AND votes.voter = ? AND "+postType+"_votes."+postType+"_id = ?", tagID, voter, postID).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func FindPaginatedPostsByTagOrderedByNumberOfVotes(tagName string, page, size int) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * size
	var posts []model.Post
	err := DB.Table("posts").
		Select("posts.*").
		Joins("JOIN post_votes ON post_votes.post_id = posts.id").
		Joins("JOIN votes ON votes.id = post_votes.vote_id").
		Joins("JOIN tags ON tags.id = votes.tag_id").
		Where("tags.name = ?", tagName).
		Group("posts.id").
		Offset(offset).
		Limit(size).
		Find(&posts).Error
	return posts, err
}

func FindPostByOwnerIdAndType(articleID uint64, postType string) (model.Post, error) {
	var post model.Post
	err := DB.Where("owner_id = ? AND owner_type = ?", articleID, postType).First(&post).Error
	return post, err
}
