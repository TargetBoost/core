package services

import (
	"core/internal/queue"
	"core/internal/repositories"
	"core/internal/services/auth"
	"core/internal/services/settings"
	"core/internal/services/storage"
	"core/internal/services/target"
	"core/internal/services/user"
)

type Services struct {
	repo *repositories.Repositories

	User     *user.Service
	Auth     *auth.Service
	Target   *target.Service
	Storage  *storage.Service
	Settings *settings.Service
}

func NewServices(repo *repositories.Repositories, lineBroker chan []queue.Task, LineAppoint chan queue.Task) *Services {
	userService := user.NewUserService(repo.User, LineAppoint)
	authService := auth.NewAuthService(repo.Auth)
	TargetService := target.NewTargetService(repo.User, repo.Feed, repo.Storage, lineBroker)
	storageService := storage.NewStorageService(repo.Storage)
	settingsService := settings.NewSettingsService(repo.Settings)

	return &Services{
		repo:     repo,
		User:     userService,
		Auth:     authService,
		Target:   TargetService,
		Storage:  storageService,
		Settings: settingsService,
	}
}
