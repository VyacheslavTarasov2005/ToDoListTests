package db

import (
	"HITS_ToDoList_Tests/internal/domain/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Task{})
}
