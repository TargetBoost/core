package settings

import (
	"core/internal/models"
	"core/internal/repositories"
)

type Service struct {
	repo *repositories.Repositories
}

func NewSettingsService(repo *repositories.Repositories) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetSettings() models.Settings {
	return s.repo.Settings.GetSettings()
}

func (s *Service) SetSettings(settings *models.Settings) {
	s.repo.Settings.SetSettings(settings)
}
