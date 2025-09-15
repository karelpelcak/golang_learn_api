package database

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"taskflow/internal/models"
)

// DB is the global database instance
var DB *gorm.DB

// ConnectDB initializes the SQLite database connection
func ConnectDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("taskflow.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate our schema
	err = DB.AutoMigrate(&models.Task{}, &models.TaskVersion{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}