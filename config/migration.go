package config

import (
	"log"

	"github.com/jinzhu/gorm"
)

func Migration() {
	databaseURL := "sqlite_test.db"

	db, err := gorm.Open("sqlite3", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&Book{})

	log.Println("Migration completed successfully")
}
