package database

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var GormConn *gorm.DB

func SetUpDB() {
	var err error = nil
	GormConn, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = GormConn.AutoMigrate(&model.Password{}, &model.User{}, &model.Auth{}, &model.Article{}, &model.Image{}, &model.ArticleCollection{}, &model.ImageCollection{}, &model.Post{})
	if err != nil {
		panic("failed to migrate database")
	}
}
