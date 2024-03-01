package database

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GormConn *gorm.DB

var dbName = "test.db"
var admin model.User

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	newDBName := os.Getenv("DB_NAME")
	if newDBName != "" {
		dbName = newDBName
	}
	if adminUsername != "" && adminPassword != "" {
		admin.Username = adminUsername
		admin.Password.ValidateAndSetPassword(adminPassword)
	}

}

func SetUpDB() {
	var err error = nil
	maxAuth := model.Auth{Level: math.MaxUint8}
	minAuth := model.Auth{Level: 0}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond * 300, // Slow SQL threshold
			LogLevel:                  logger.Error,           // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,                   // Include params in the SQL log
			Colorful:                  false,                  // Disable color
		},
	)

	GormConn, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	err = GormConn.AutoMigrate(&model.Password{}, &model.User{}, &model.Auth{}, &model.Article{},
		&model.Image{}, &model.ArticleCollection{}, &model.ImageCollection{}, &model.Post{},
		&model.Tag{}, &model.Section{}, &model.Profile{})
	if err != nil {
		panic("failed to migrate database")
	}
	//Save auths
	err = GormConn.Model(&model.Auth{}).Where("level = ?", maxAuth.Level).First(&maxAuth).Error
	if err != nil {
		err = GormConn.Create(&maxAuth).Error
		if err != nil {
			log.Fatalln("Error creating maxAuth: ", err)
		}
	}
	err = GormConn.Model(&model.Auth{}).Where("level = ?", minAuth.Level).First(&minAuth).Error
	if err != nil {
		err = GormConn.Create(&minAuth).Error
		if err != nil {
			log.Fatalln("Error creating maxAuth: ", err)
		}
	}
	//Save admin
	var adminInDB model.User
	err = GormConn.Where("username = ?", admin.Username).First(&adminInDB).Error
	if err != nil {
		adminInDB = admin
		err = GormConn.Create(&adminInDB).Error
		if err != nil {
			log.Fatalln("Error creating admin: ", err)
		}
	}
	adminInDB.Password.Hash = admin.Password.Hash
	adminInDB.Auth = maxAuth
	err = GormConn.Save(&adminInDB).Error
	if err != nil {
		log.Fatalln("Error saving admin: ", err)
	}

}
