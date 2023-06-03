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

func (s *Service) SetChatMembers(cid int64, title, userName, photoLink, bio string) {
	s.storageRepository.SetChatMembers(cid, title, userName, photoLink, bio)
}
