package services

import (
	"core/internal/queue"
	"core/internal/repositories"
	"core/internal/services/settings"
	"core/internal/services/storage"
	"core/internal/services/target"
	"core/internal/services/user"
	"core/internal/tg/bot"
)

type Services struct {
	repo *repositories.Repositories

	User     *user.Service
	Target   *target.Service
	Storage  *storage.Service
	Settings *settings.Service
}

func NewServices(repo *repositories.Repositories, lineBroker chan []queue.Task, lineAppoint chan queue.Task, trackMessages chan bot.Message) *Services {
	userService := user.NewUserService(repo, lineAppoint, trackMessages)
	TargetService := target.NewTargetService(repo, lineBroker, trackMessages)
	storageService := storage.NewStorageService(repo)
	settingsService := settings.NewSettingsService(repo)

	return &Services{
		repo:     repo,
		User:     userService,
		Target:   TargetService,
		Storage:  storageService,
		Settings: settingsService,
	}
}
