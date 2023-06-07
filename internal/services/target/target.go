package target

import (
	"core/internal/models"
	"core/internal/queue"
	"core/internal/repositories"
	"core/internal/tg/bot"
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
	repo          *repositories.Repositories
	lineBroker    chan []queue.Task
	trackMessages chan bot.Message
}

func NewTargetService(repo *repositories.Repositories, lineBroker chan []queue.Task, trackMessages chan bot.Message) *Service {
	return &Service{
		repo:          repo,
		lineBroker:    lineBroker,
		trackMessages: trackMessages,
	}
}

func (s *Service) GetTargets(uid uint) []models.TargetService {
	targets := func(t []models.Target, f func(t models.Target) models.TargetService) []models.TargetService {
		result := make([]models.TargetService, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.repo.Target.GetTargets(uid), models.MapToTarget)

	return targets
}

func (s *Service) GetTarget(tid uint) models.TargetService {
	t := s.repo.Target.GetTargetByID(tid)

	return models.TargetService{
		ID:         t.ID,
		UID:        t.UID,
		Title:      t.Title,
		Link:       t.Link,
		Icon:       t.Icon,
		Status:     t.Status,
		Count:      strconv.Itoa(int(t.Count)),
		Total:      strconv.Itoa(int(t.Total)),
		Cost:       t.Cost,
		Cause:      t.Cause,
		TotalPrice: strconv.Itoa(int(t.TotalPrice)),
	}
}

func (s *Service) GetTaskByID(id uint) models.QueueToService {
	t := s.repo.Queue.GetTaskByID(int64(id))

	return models.QueueToService{
		Status: t.Status,
	}
}

func (s *Service) GetTargetsToAdmin() []models.TargetService {
	targets := func(t []models.TargetToAdmin, f func(t models.TargetToAdmin) models.TargetService) []models.TargetService {
		result := make([]models.TargetService, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.repo.Target.GetTargetsToAdmin(), models.MapToTargetAdmin)

	return targets
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

func (s *Service) UpdateTarget(id uint, status int64) {
	t := s.repo.Target.GetTargetByID(id)
	t.Status = status

	s.repo.Target.UpdateTarget(id, &t)

	if status == 1 {
		u := s.repo.Account.GetAllUsers()
		for _, v := range u {
			st := strings.ToLower(strings.Split(v.Tg, "@")[len(strings.Split(v.Tg, "@"))-1])
			cm := s.repo.Queue.GetChatMembersByUserName(st)
			m := bot.Message{
				Type: 100,
				CID:  cm.CID,
			}

			if v.Execute {
				s.trackMessages <- m
			}
		}
	}
}

func (s *Service) GetChatID(id uint) (int64, float64) {
	tu := s.GetTarget(id)
	st := strings.Split(tu.Link, "/")[len(strings.Split(tu.Link, "/"))-1]

	ch := s.repo.Queue.GetChatMembersByUserName(st)
	return ch.CID, tu.Cost
}

func (s *Service) GetUserID(id uint) int64 {
	tu := s.repo.Account.GetUserByID(int64(id))
	st := strings.Split(tu.Tg, "@")[len(strings.Split(tu.Tg, "@"))-1]

	ch := s.repo.Queue.GetChatMembersByUserName(st)
	return ch.CID
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
	st = strings.ToLower(st)
	ch := s.repo.Queue.GetChatMembersByUserName(st)
	if ch.CID == 0 {
		return errors.New("Вы не добавили нашего бота в этот телеграм канал")
	}

	u := s.repo.Account.GetUserByID(int64(UID))
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

	s.repo.Account.UpdateUserBalance(u.ID, u.Balance)

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

	tt := s.repo.Target.CreateTarget(&t)

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
