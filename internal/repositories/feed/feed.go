package feed

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewFeedRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
