package database

import (
	"log"
	"taskflow/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

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
