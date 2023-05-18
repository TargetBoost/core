package target

import (
	"core/internal/models"
	"core/internal/queue"
	"core/internal/repositories/target"
	"github.com/ivahaev/go-logger"
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
	lineBroker       chan []queue.Task
}

func NewTargetService(feedRepository *target.Repository, lineBroker chan []queue.Task) *Service {
	return &Service{
		TargetRepository: feedRepository,
		lineBroker:       lineBroker,
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

func (s *Service) GetTargetsToExecutor() []models.TargetService {
	targets := func(t []models.Target, f func(t models.Target) models.TargetService) []models.TargetService {
		result := make([]models.TargetService, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.TargetRepository.GetTargetsToExecutor(), models.MapToTarget)

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
		Status: 0,
		Count:  0,
		Total:  target.Total,
		Cost:   target.Cost,
		Type:   target.Type,
	}

	tt := s.TargetRepository.CreateTarget(&t)

	var q []queue.Task

	logger.Info(tt)

	var i int64 = 0
	for i = 0; i < t.Total; i++ {
		q = append(q, queue.Task{
			ID:     tt.ID,
			Cost:   t.Cost,
			Title:  t.Title,
			Status: t.Status,
		})
	}

	s.lineBroker <- q
}
