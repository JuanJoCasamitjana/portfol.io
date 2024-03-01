package database

import (
	"errors"

	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"gorm.io/gorm"
)

var ErrIDInvalid = errors.New("ID is invalid")

func GetProfileById(id uint64) (model.Profile, error) {
	var profile model.Profile
	result := GormConn.First(&profile, id)
	return profile, result.Error
}

func GetProfileByUserID(id uint64) (model.Profile, error) {
	var profile model.Profile
	result := GormConn.Where("user_id = ?", id).First(&profile)
	return profile, result.Error
}

func CreateProfile(profile *model.Profile) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(profile)
		return result.Error
	})
}

func UpdateProfile(profile *model.Profile) error {
	if profile.ID == 0 {
		return ErrIDInvalid
	}
	return GormConn.Transaction(func(tx *gorm.DB) error {
		result := tx.Save(profile)
		return result.Error
	})
}

func DeleteProfileByID(id uint64) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&model.Profile{}, id)
		return result.Error
	})
}

func DeleteProfileByUserID(id uint64) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("user_id = ?", id).Delete(&model.Profile{})
		return result.Error
	})
}

func DeleteProfile(profile *model.Profile) error {
	return GormConn.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(profile)
		return result.Error
	})
}
