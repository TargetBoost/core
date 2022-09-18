package repositories

import (
	"core/internal/models"
	"core/internal/repositories/auth"
	"core/internal/repositories/feed"
	"core/internal/repositories/storage"
	"core/internal/repositories/user"
	"github.com/ivahaev/go-logger"
	"gorm.io/gorm"
)

type Repositories struct {
	db *gorm.DB

	User    *user.Repository
	Auth    *auth.Repository
	Feed    *feed.Repository
	Storage *storage.Repository
}

func NewRepositories(db *gorm.DB) *Repositories {
	err := db.AutoMigrate(
		&models.User{},
		&models.Feed{},
		&models.FileStorage{},
	)
	if err != nil {
		logger.Error(err)
	}

	userRepository := user.NewUserRepository(db)
	authRepository := auth.NewAuthRepository(db)
	feedRepository := feed.NewFeedRepository(db)
	storageRepository := storage.NewStorageRepository(db)

	return &Repositories{
		db:      db,
		User:    userRepository,
		Auth:    authRepository,
		Feed:    feedRepository,
		Storage: storageRepository,
	}
}
