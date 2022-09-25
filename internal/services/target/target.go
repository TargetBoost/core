package target

import (
	"core/internal/models"
	"core/internal/repositories/target"
)

type Service struct {
	TargetRepository *target.Repository
}

func NewTargetService(feedRepository *target.Repository) *Service {
	return &Service{
		TargetRepository: feedRepository,
	}
}

func (s *Service) GetTargets(uid uint) []models.TargetService {
	targets := func(t []models.Target, f func(t models.Target) models.TargetService) []models.TargetService {
		result := make([]models.TargetService, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.TargetRepository.GetTargets(uid), models.MapToTarget)

	return targets
}

func (s *Service) CreateTarget(UID uint, target *models.TargetService) {

	t := models.Target{
		UID:    UID,
		Title:  target.Title,
		Link:   target.Link,
		Icon:   target.Icon,
		Status: "check",
		Count:  target.Count,
		Cost:   target.Cost,
		Type:   target.Type,
	}

	s.TargetRepository.CreateTarget(&t)
}
