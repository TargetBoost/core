package services

import (
	"core/internal/models"
	"core/internal/repositories"
	"core/internal/services/account"
	"core/internal/services/blog"
	"core/internal/services/queue"
	"core/internal/services/settings"
	"core/internal/services/storage"
	"core/internal/services/target"
	"core/internal/target_broker"
	"core/internal/transport/tg/bot"
)

type Account interface {
	CreateUser(user models.CreateUser) (*models.User, error)
	IsAuth(token string) (uint, bool)
	IsAdmin(token string) bool
	UpdateUserBalance(id int64, cost float64)
	UpdateUser(uid uint, b float64)
	GetAllUsers() []models.UserService
	GetTasksCashesUser(uid uint) []models.TaskCashToService
	GetTasksCashesAdmin() []models.TaskCashToService
	GetUserByID(id int64) models.UserService
	GetUserByToken(token string) models.UserService
	AuthUser(user models.AuthUser) (*models.User, error)
	CreateTaskCashes(uid int64, task models.TaskCashToUser) error
	UpdateTaskCashes(task models.TaskCashToService)
	CreateTransaction(t *models.TransactionToService)
	UpdateTransaction(t *models.TransactionToService)
	GetTransaction(build string) *models.TransactionToService
}

type Target interface {
	GetTargets(uid uint) []models.TargetService
	GetTarget(tid uint) models.TargetService
	GetTargetsToAdmin() []models.TargetService
	UpdateTarget(id uint, status int64)
	GetUserID(id uint) int64
	CreateTarget(UID uint, target *models.TargetService) error
	GetProfit() float64
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

type Blog interface {
	GetBlog() []models.Blog
	AddComment(c models.Comment)
}

type Services struct {
	Account  Account
	Queue    Queue
	Target   Target
	Storage  Storage
	Settings Settings
	Blog     Blog
}

func NewServices(repo *repositories.Repositories, lineBroker chan []target_broker.Task, lineAppoint chan target_broker.Task, trackMessages chan bot.Message) *Services {
	accountService := account.NewAccountService(repo, lineAppoint, trackMessages)
	targetService := target.NewTargetService(repo, lineBroker, trackMessages)
	storageService := storage.NewStorageService(repo)
	settingsService := settings.NewSettingsService(repo)
	queueService := queue.NewQueueService(repo)
	blogService := blog.NewBlogService(repo)

	return &Services{
		Account:  accountService,
		Queue:    queueService,
		Target:   targetService,
		Storage:  storageService,
		Settings: settingsService,
		Blog:     blogService,
	}
}
