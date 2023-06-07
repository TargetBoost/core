package services

import (
	"core/internal/models"
	"core/internal/repositories"
	"core/internal/services/account"
	"core/internal/services/queue"
	"core/internal/services/settings"
	"core/internal/services/storage"
	"core/internal/services/target"
	"core/internal/target_broker"
	"core/internal/transport/tg/bot"
)

type Account interface {
}

type Target interface {
	GetTargets(uid uint) []models.TargetService
	GetTarget(tid uint) models.TargetService
	GetTargetsToAdmin() []models.TargetService
	UpdateTarget(id uint, status int64)
	GetUserID(id uint) int64
	CreateTarget(UID uint, target *models.TargetService) error
}

type Storage interface {
	GetFileByKey(key string) *models.FileStorage
	SetChatMembers(cid, count int64, title, userName, photoLink, bio string)
	CallBackVK(code, token string) error
}

type Queue interface {
	GetTaskByID(id uint) models.QueueToService
	GetTargetsToExecutor(uid int64) []models.QueueToService
	UpdateTaskStatus(id uint)
	GetChatID(id uint) (int64, float64)
}

type Settings interface {
	GetSettings() models.Settings
	SetSettings(settings *models.Settings)
}

type Services struct {
	Account  *account.Service
	Queue    Queue
	Target   Target
	Storage  Storage
	Settings Settings
}

func NewServices(repo *repositories.Repositories, lineBroker chan []target_broker.Task, lineAppoint chan target_broker.Task, trackMessages chan bot.Message) *Services {
	accountService := account.NewAccountService(repo, lineAppoint, trackMessages)
	targetService := target.NewTargetService(repo, lineBroker, trackMessages)
	storageService := storage.NewStorageService(repo)
	settingsService := settings.NewSettingsService(repo)
	queueService := queue.NewQueueService(repo)

	return &Services{
		Account:  accountService,
		Queue:    queueService,
		Target:   targetService,
		Storage:  storageService,
		Settings: settingsService,
	}
}
