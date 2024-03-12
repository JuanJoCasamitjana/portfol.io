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
