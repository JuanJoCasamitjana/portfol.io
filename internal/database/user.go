package database

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"gorm.io/gorm"
)

func SaveUser(user *model.User) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := GormConn.Preload("Password").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint64) (*model.User, error) {
	var user model.User
	err := GormConn.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SaveAuth(auth *model.Auth) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(auth).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetAuthByID(id uint64) (*model.Auth, error) {
	var auth model.Auth
	err := GormConn.Where("id = ?", id).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func UpdateUserById(id uint64, upadatedFields map[string]interface{}) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		var user model.User
		err := tx.Where("id = ?", user.ID).First(&user).Error
		if err != nil {
			return err
		}
		if err := tx.Model(&user).Updates(upadatedFields).Error; err != nil {
			return err
		}
		return nil
	})
}

func DeleteUser(user *model.User) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		if result := tx.Delete(user); result.Error != nil {
			return result.Error
		}
		return nil
	})
}
