package repositories

import (
	"core/internal/models"
	"core/internal/repositories/account"
	"core/internal/repositories/queue"
	"core/internal/repositories/settings"
	"core/internal/repositories/storage"
	"core/internal/repositories/target"

	"github.com/ivahaev/go-logger"
	"gorm.io/gorm"
)

type Account interface {
	GetAllUsers() []models.User
	GetUserByID(id int64) models.User
	GetUserByToken(token string) models.User
	GetUserByTgAndPassword(tg, pass string) models.User
	IsUserByTg(tg string) bool
	UpdateUser(user models.User)
	UpdateUserBalance(uid uint, balance float64)
	CreateUser(user *models.CreateUser) error
	CreateTaskCache(task models.TaskCash)
	UpdateTaskCache(task models.TaskCash)
	GetTaskCacheByUID(uid uint) []models.TaskCash
	GetTaskCacheByID(id uint) models.TaskCash
	GetTaskCacheToAdmin() []models.TaskCash
	CreateTransaction(t *models.TransactionToService)
	UpdateTransaction(t *models.TransactionToService)
	GetTransaction(build string) models.Transaction
	IsAuth(token string) (uint, bool)
}

type Target interface {
	CreateTarget(target *models.Target) *models.Target
	GetTargets(uid uint) []models.Target
	GetTargetByID(uid uint) models.Target
	GetTargetsToAdmin() []models.TargetToAdmin
	UpdateTarget(id uint, target *models.Target)
}

type Queue interface {
	GetTargetsToExecutor(uid int64) []models.Queue
	GetTaskByID(id int64) models.Queue
	CreateTask(queue *models.Queue)
	UpdateTaskStatus(q models.Queue)
	UpdateTask(q models.Queue)
	GetUniqueTask() []models.Queue
	GetTaskDISTINCTInWork() []models.Queue
	GetTaskForUserUID(uid uint, tid uint) []models.Queue
	GetTaskDISTINCTIsWorkForUser(uid int64) []models.QueueToExecutors
	GetChatMembersByUserName(userName string) models.ChatMembersChanel
}

type Storage interface {
	GetFileByKey(key string) *models.FileStorage
	SetChatMembers(cid, countSub int64, title, userName, photoLink, bio string)
	GetStatisticTargetsOnExecutesIsTrue() []models.StatisticTargetsOnExecutesIsTrue
}

type Settings interface {
	GetSettings() models.Settings
	SetSettings(s *models.Settings)
}

type Repositories struct {
	Account  Account
	Target   Target
	Queue    Queue
	Storage  Storage
	Settings Settings
}

func NewRepositories(db *gorm.DB) *Repositories {
	err := db.AutoMigrate(
		&models.User{},
		&models.Target{},
		&models.FileStorage{},
		&models.Settings{},
		&models.TargetToExecutors{},
		&models.Queue{},
		&models.ChatMembersChanel{},
		&models.TaskCash{},
		&models.Transaction{},
	)
	if err != nil {
		logger.Error(err)
	}

	accountRepository := account.NewAccountRepository(db)
	targetRepository := target.NewTargetRepository(db)
	queueRepository := queue.NewQueueRepository(db)
	storageRepository := storage.NewStorageRepository(db)
	settingsRepository := settings.NewSettingsRepository(db)

	return &Repositories{
		Account:  accountRepository,
		Target:   targetRepository,
		Queue:    queueRepository,
		Storage:  storageRepository,
		Settings: settingsRepository,
	}
}
