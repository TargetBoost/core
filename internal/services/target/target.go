package target

import (
	"core/internal/models"
	"core/internal/repositories/target"
)

const (
	vkCommunity  = "vk_community"
	vkLike       = "vk_like"
	vkAddFriends = "vk_add_friends"
	tgCommunity  = "tg_community"
	ytChanel     = "yt_chanel"
	ytWatch      = "yt_watch"
	ytLike       = "yt_like"
	ytDislike    = "yt_dislike"
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

func (s *Service) GetTargetsToAdmin() []models.TargetService {
	targets := func(t []models.Target, f func(t models.Target) models.TargetService) []models.TargetService {
		result := make([]models.TargetService, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.TargetRepository.GetTargetsToAdmin(), models.MapToTarget)

	return targets
}

func (s *Service) GetTargetsToExecutor() []models.TargetServiceToExecutors {
	targets := func(t []models.TargetToExecutors, f func(t models.TargetToExecutors) models.TargetServiceToExecutors) []models.TargetServiceToExecutors {
		result := make([]models.TargetServiceToExecutors, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.TargetRepository.GetTargetsToExecutor(), models.MapToTargetExecutors)

	return targets
}

func (s *Service) CreateTarget(UID uint, target *models.TargetService) {
	var title string
	switch target.Type {
	case vkCommunity:
		title = "Вступить в сообщество"
		break
	case vkLike:
		title = "Поставить лайк на запись"
		break
	case vkAddFriends:
		title = "Добавить в друзья"
		break
	case tgCommunity:
		title = "Подписаться на канал"
		break
	case ytChanel:
		title = "Подписаться на канал"
		break
	case ytWatch:
		title = "Посмотреть видео"
		break
	case ytLike:
		title = "Поставить лайк"
		break
	case ytDislike:
		title = "Поставить дизлайк"
		break
	}

	t := models.Target{
		UID:    UID,
		Title:  title,
		Link:   target.Link,
		Icon:   target.Icon,
		Status: "check",
		Count:  0,
		Total:  target.Total,
		Cost:   target.Cost,
		Type:   target.Type,
	}

	s.TargetRepository.CreateTarget(&t)
}
