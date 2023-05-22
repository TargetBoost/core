package target

import (
	"core/internal/models"
	"core/internal/queue"
	"core/internal/repositories/target"
	"core/internal/repositories/user"
	"errors"
	"github.com/ivahaev/go-logger"
	"strconv"
	"strings"
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
	UserRepository   *user.Repository
	lineBroker       chan []queue.Task
}

func NewTargetService(userRepository *user.Repository, feedRepository *target.Repository, lineBroker chan []queue.Task) *Service {
	return &Service{
		TargetRepository: feedRepository,
		UserRepository:   userRepository,
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
	targets := func(t []models.TargetToAdmin, f func(t models.TargetToAdmin) models.TargetService) []models.TargetService {
		result := make([]models.TargetService, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.TargetRepository.GetTargetsToAdmin(), models.MapToTargetAdmin)

	return targets
}

func (s *Service) GetTargetsToExecutor(uid int64) []models.QueueToService {
	targets := func(t []models.QueueToExecutors, f func(t models.QueueToExecutors) models.QueueToService) []models.QueueToService {
		result := make([]models.QueueToService, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.TargetRepository.GetTaskDISTINCTIsWorkForUser(uid), models.MapToQueueExecutors)

	return targets
}

func (s *Service) UpdateTarget(id uint, status int64) {
	t := s.TargetRepository.GetTargetByID(id)
	t.Status = status

	s.TargetRepository.UpdateTarget(id, &t)
}

func (s *Service) CreateTarget(UID uint, target *models.TargetService) error {
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

	st := strings.Split(target.Link, "/")[len(strings.Split(target.Link, "/"))-1]
	ch := s.TargetRepository.GetChatMembersByUserName(st)
	if ch.CID == 0 {
		return errors.New("Вы не добавили нашего бота в этот телеграм канал")
	}

	u := s.UserRepository.GetUserByID(int64(UID))
	if u.ID == 0 {
		return errors.New("user not found")
	}

	if target.Cost < 0 {
		return errors.New("error create task")
	}

	if u.Balance < target.Cost {
		return errors.New("Вашего баланса недостаточно для создания рекламной кампании")
	}

	ft, err := strconv.Atoi(target.Total)
	if err != nil {
		return err
	}

	tl := target.Cost * float64(ft)

	u.Balance = u.Balance - tl
	if u.Balance < 0 {
		return errors.New("Вашего баланса недостаточно для создания рекламной кампании")
	}

	s.UserRepository.UpdateUser(u)

	t := models.Target{
		UID:        UID,
		Title:      title,
		Link:       target.Link,
		Icon:       target.Icon,
		Status:     0,
		Count:      0,
		Total:      float64(ft),
		Cost:       target.Cost,
		Type:       target.Type,
		TotalPrice: tl,
	}

	tt := s.TargetRepository.CreateTarget(&t)

	var q []queue.Task

	logger.Info(tt)

	var i float64 = 0
	for i = 0; i < t.Total; i++ {
		q = append(q, queue.Task{
			TID:    tt.ID,
			Cost:   t.Cost,
			Title:  t.Title,
			Status: t.Status,
		})
	}

	s.lineBroker <- q
	return nil
}
