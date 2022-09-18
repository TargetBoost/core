package storage

import (
	"core/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewStorageRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetFileByKey(key string) *models.FileStorage {
	var fileStorage models.FileStorage

	r.db.Table("file_storages").Where("key = ?", key).Find(&fileStorage)

	if fileStorage.Key == "" {
		return nil
	}

	return &fileStorage
}
