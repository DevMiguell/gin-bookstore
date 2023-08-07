package config

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func getTestDatabasePath() string {
	return filepath.Join("../..", "sqlite_test.db")
}

func ConnectDatabase() {
	databaseURL := "sqlite.db"

	if IsTestEnvironment() {
		databaseURL = getTestDatabasePath()
	}

	var err error
	database, err := gorm.Open("sqlite3", databaseURL)

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&User{})
	database.AutoMigrate(&Book{})

	DB = database
}

func IsTestEnvironment() bool {
	return os.Getenv("ENVIRONMENT") == "test"
}
