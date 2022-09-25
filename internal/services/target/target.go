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
	var title string
	switch target.Type {
	case "vk_community":
		title = "Вступить в сообщество"
		break
	case "vk_like":
		title = "Поставить лайк на запись"
		break
	case "vk_add_friends":
		title = "Добавить в друзья"
		break
	}

	t := models.Target{
		UID:    UID,
		Title:  title,
		Link:   target.Link,
		Icon:   target.Icon,
		Status: "check",
		Count:  target.Count,
		Cost:   target.Cost,
		Type:   target.Type,
	}

	s.TargetRepository.CreateTarget(&t)
}
