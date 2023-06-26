package queue

import (
	"core/internal/models"
	"core/internal/repositories"
	"strings"
)

type Service struct {
	repo *repositories.Repositories
}

func (s *Service) GetTaskByID(id uint) models.QueueToService {
	t := s.repo.Queue.GetTaskByID(int64(id))

	return models.QueueToService{
		Status: t.Status,
	}
}

func (s *Service) GetTargetsToExecutor(uid int64) []models.QueueToService {
	us := s.repo.Account.GetUserByID(uid)
	st := strings.ToLower(strings.Split(us.Tg, "@")[len(strings.Split(us.Tg, "@"))-1])

	stu := s.repo.Queue.GetChatMembersByUserName(st)

	if stu.CID == 0 {
		return []models.QueueToService{}
	}

	targets := func(t []models.QueueToExecutors, f func(t models.QueueToExecutors) models.QueueToService) []models.QueueToService {
		result := make([]models.QueueToService, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.repo.Queue.GetTaskDISTINCTIsWorkForUser(uid), models.MapToQueueExecutors)

	return targets
}

func (s *Service) UpdateTaskStatus(id uint) {
	var q models.Queue
	q.ID = id
	q.Status = 3
	s.repo.Queue.UpdateTaskStatus(q)
}

func (s *Service) GetChatID(id uint) (int64, float64) {
	tu := s.repo.Target.GetTargetByID(id)
	st := strings.Split(tu.Link, "/")[len(strings.Split(tu.Link, "/"))-1]

	ch := s.repo.Queue.GetChatMembersByUserName(st)
	return ch.CID, 0
}

func NewQueueService(repo *repositories.Repositories) *Service {
	return &Service{
		repo: repo,
	}
}
