package database

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"gorm.io/gorm"
)

func FindUserById(id uint64) (model.User, error) {
	var user model.User
	result := DB.Where("id = ?", id).First(&user)
	return user, result.Error
}

func CreateUser(user *model.User) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(user)
		return result.Error
	})
}

func FindUserByUsername(username string) (model.User, error) {
	var user model.User
	result := DB.Where("username = ?", username).First(&user)
	return user, result.Error
}

func UpdateUser(user *model.User) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Save(user)
		return result.Error
	})
}

func FindUserByEmail(email string) (model.User, error) {
	var user model.User
	result := DB.Where("email = ?", email).First(&user)
	return user, result.Error
}

func FindSectionsByUser(username string) ([]model.Section, error) {
	var sections []model.Section
	result := DB.Where("owner = ?", username).Find(&sections)
	return sections, result.Error
}

func FindPostsByUserPaginated(username string, page, page_size int) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * page_size
	var posts []model.Post
	result := DB.Where("author = ? AND published = ?", username, true).Order("updated_at desc").Offset(offset).
		Limit(page_size).Find(&posts).Error
	return posts, result
}

func FindSectionByUsernameAndName(username, name string) (model.Section, error) {
	var section model.Section
	result := DB.Where("owner = ? AND name = ?", username, name).Preload("Posts").First(&section)
	return section, result.Error
}

func FindPostsByUserAndSectionPaginated(username, section string, page, page_size int) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * page_size
	var posts []model.Post
	//Section has a many to many relationship with posts and posts has no foreign key to section
	//So we need to do a subquery to get the posts
	result := DB.Where("author = ? AND published = true AND id IN (SELECT post_id FROM section_posts WHERE section_id = (SELECT id FROM sections WHERE owner = ? AND name = ?))", username, username, section).
		Order("updated_at desc").Offset(offset).Limit(page_size).Find(&posts).Error
	return posts, result
}

func CreateSection(section *model.Section) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(section)
		return result.Error
	})
}

func AddPostToSection(section *model.Section, post *model.Post) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(section).Association("Posts").Append(post)
	})
}

func RemovePostFromSection(section *model.Section, post *model.Post) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(section).Association("Posts").Delete(post)
	})
}

func FindPostsByUserNotInSectionPaginated(username, section string, page, page_size int) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * page_size
	var posts []model.Post
	//Section has a many to many relationship with posts and posts has no foreign key to section
	//So we need to do a subquery to get the posts
	result := DB.Where("author = ? AND published = true AND id NOT IN (SELECT post_id FROM section_posts WHERE section_id = (SELECT id FROM sections WHERE owner = ? AND name = ?))", username, username, section).
		Order("updated_at desc").Offset(offset).Limit(page_size).Find(&posts).Error
	return posts, result
}

func DeleteSectionByUsernameAndName(username, name string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("owner = ? AND name = ?", username, name).Delete(&model.Section{})
		return result.Error
	})
}

func DeleteUser(user *model.User) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("author = ?", user.Username).Delete(&model.Post{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("owner = ?", user.Username).Delete(&model.Section{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("author = ?", user.Username).Delete(&model.Article{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("author = ?", user.Username).Delete(&model.Gallery{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("owner = ?", user.Username).Delete(&model.Image{}).Error
		if err != nil {
			return err
		}
		err = tx.Where("author = ?", user.Username).Delete(&model.Project{}).Error
		if err != nil {
			return err
		}
		return tx.Delete(user).Error

	})
}

func FollowUser(follower_follow_list *model.FollowList, followed *model.User) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(follower_follow_list).Where("owner = ?", follower_follow_list.Owner).
			Association("Following").Append(followed)
	})
}

func UnfollowUser(follower_follow_list *model.FollowList, followed *model.User) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(follower_follow_list).Association("Following").Delete(followed)
	})
}

func FindFollowingPostsPaginated(user model.User, page, pageSize int) ([]model.Post, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	var posts []model.Post
	result := DB.Where("author IN (SELECT username FROM follows WHERE owner = ?) AND published = true", user.Username).
		Order("updated_at desc").Offset(offset).Limit(pageSize).Find(&posts).Error
	return posts, result
}

func FindFollowListByUsername(username string) (model.FollowList, error) {
	var followList model.FollowList
	result := DB.Where("owner = ?", username).Preload("Following").First(&followList)
	return followList, result.Error
}

func CreateFollowList(followList *model.FollowList) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(followList)
		return result.Error
	})
}

func FindUsersPaginated(page, pageSize int) ([]model.User, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	var users []model.User
	result := DB.Offset(offset).Limit(pageSize).Find(&users)
	return users, result.Error
}

func CountUsers() (int64, error) {
	var count int64
	result := DB.Model(&model.User{}).Count(&count)
	return count, result.Error
}

func FindUsersPaginatedBySearch(search string, page, pageSize int) ([]model.User, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	var users []model.User
	result := DB.Where("username LIKE ?", "%"+search+"%").Offset(offset).Limit(pageSize).Find(&users)
	return users, result.Error
}
