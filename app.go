package main

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var KEY string
var db *gorm.DB

func init() {
	KEY = os.Getenv("JWT_KEY")
	if KEY == "" {
		log.Fatalf("The JWT secret key is not set.")
	}
	var err error

	// ? bd should be a more generalised storage interface
	db, err = gorm.Open(sqlite.Open("onesgo.db"), &gorm.Config{
		// ? disable this for product deployment
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	db.AutoMigrate(&User{})
}

func main() {
	router := setupRouter()
	router.Run(":8080")
}
