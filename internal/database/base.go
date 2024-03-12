package database

import (
	"log"
	"os"
	"time"

	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbName = "dev.db"
var DB *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	newDBName := os.Getenv("DB_NAME")
	if newDBName != "" {
		dbName = newDBName
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond * 300, // Slow SQL threshold
			LogLevel:                  logger.Error,           // Log level
			IgnoreRecordNotFoundError: false,                  // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,                  // Include params in the SQL log
			Colorful:                  false,                  // Disable color
		},
	)
	DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = DB.AutoMigrate(model.User{}, model.Article{}, model.Project{}, model.Image{}, model.Gallery{}, model.Post{})
	if err != nil {
		log.Fatal(err)
	}
}
