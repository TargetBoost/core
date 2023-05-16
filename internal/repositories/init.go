package repositories

import (
	"core/internal/models"
	"core/internal/repositories/auth"
	"core/internal/repositories/settings"
	"core/internal/repositories/storage"
	"core/internal/repositories/target"
	"core/internal/repositories/user"
	"github.com/ivahaev/go-logger"
	"gorm.io/gorm"
)

type Repositories struct {
	db *gorm.DB

	User     *user.Repository
	Auth     *auth.Repository
	Feed     *target.Repository
	Storage  *storage.Repository
	Settings *settings.Repository
}

func NewRepositories(db *gorm.DB) *Repositories {
	err := db.AutoMigrate(
		&models.User{},
		&models.Target{},
		&models.FileStorage{},
		&models.Settings{},
		&models.TargetToExecutors{},
	)
	if err != nil {
		logger.Error(err)
	}

	userRepository := user.NewUserRepository(db)
	authRepository := auth.NewAuthRepository(db)
	feedRepository := target.NewTargetRepository(db)
	storageRepository := storage.NewStorageRepository(db)
	settingsRepository := settings.NewSettingsRepository(db)

	return &Repositories{
		db:       db,
		User:     userRepository,
		Auth:     authRepository,
		Feed:     feedRepository,
		Storage:  storageRepository,
		Settings: settingsRepository,
	}
}
