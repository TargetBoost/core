package queue

import (
	"core/internal/repositories"
)

type Service struct {
	repository *repositories.Repositories
}

func NewQueueService(repository *repositories.Repositories) *Service {
	return &Service{
		repository: repository,
	}
}
