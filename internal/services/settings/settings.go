package settings

import (
	"core/internal/models"
	"core/internal/repositories/settings"
)

type Service struct {
	settingsRepository *settings.Repository
}

func NewSettingsService(settingsRepository *settings.Repository) *Service {
	return &Service{
		settingsRepository: settingsRepository,
	}
}

func (s *Service) GetSettings() models.Settings {
	return s.settingsRepository.GetSettings()
}

func (s *Service) SetSettings(settings *models.Settings) {
	s.settingsRepository.SetSettings(settings)
}
