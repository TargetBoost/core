package settings

import (
	"core/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetSettings() models.Settings {
	var s models.Settings
	r.db.Table("settings").Find(&s)

	return s
}

func (r *Repository) SetSettings(s *models.Settings) {
	r.db.Table("settings").UpdateColumns(&s)
}
