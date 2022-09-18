package storage

import (
	"core/internal/models"
	"core/internal/repositories/storage"
)

type Service struct {
	storageRepository *storage.Repository
}

func NewStorageService(storageRepository *storage.Repository) *Service {
	return &Service{
		storageRepository: storageRepository,
	}
}

func (s *Service) GetFileByKey(key string) *models.FileStorage {
	return s.storageRepository.GetFileByKey(key)
}
