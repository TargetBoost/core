package services

import (
	"core/internal/repositories"
	"core/internal/services/auth"
	"core/internal/services/feed"
	"core/internal/services/settings"
	"core/internal/services/storage"
	"core/internal/services/user"
)

type Services struct {
	repo *repositories.Repositories

	User     *user.Service
	Auth     *auth.Service
	Feed     *feed.Service
	Storage  *storage.Service
	Settings *settings.Service
}

func NewServices(repo *repositories.Repositories) *Services {
	userService := user.NewUserService(repo.User)
	authService := auth.NewAuthService(repo.Auth)
	feedService := feed.NewFeedService(repo.Feed)
	storageService := storage.NewStorageService(repo.Storage)
	settingsService := settings.NewSettingsService(repo.Settings)

	return &Services{
		repo:     repo,
		User:     userService,
		Auth:     authService,
		Feed:     feedService,
		Storage:  storageService,
		Settings: settingsService,
	}
}
