package database

import (
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
}
