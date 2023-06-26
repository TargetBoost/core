package target

import (
	"core/internal/models"
	"core/internal/repositories"
	"core/internal/target_broker"
	"core/internal/transport/tg/bot"
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
	lineBroker    chan []target_broker.Task
	trackMessages chan bot.Message
}

func (s *Service) GetProfit() float64 {
	return s.repo.Target.GetProfit()
}

func NewTargetService(repo *repositories.Repositories, lineBroker chan []target_broker.Task, trackMessages chan bot.Message) *Service {
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
		NameCompany:        t.NameCompany,
		DescriptionCompany: t.DescriptionCompany,
		Type:               t.Type,
		Link:               t.Link,
		Limit:              t.Limit,
		TypeAd:             t.TypeAd,
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

func (s *Service) GetUserID(id uint) int64 {
	tu := s.repo.Account.GetUserByID(int64(id))
	st := strings.Split(tu.Tg, "@")[len(strings.Split(tu.Tg, "@"))-1]

	ch := s.repo.Queue.GetChatMembersByUserName(st)
	return ch.CID
}

func (s *Service) CreateTarget(UID uint, target *models.TargetService) error {

	t := models.Target{
		UID:         UID,
		NameCompany: target.NameCompany,
		Link:        target.Link,
		TypeAd:      target.TypeAd,
		Type:        target.Type,
	}

	_ = s.repo.Target.CreateTarget(&t)

	//var q []target_broker.Task
	//
	//logger.Info(tt)

	//var i float64 = 0
	//for i = 0; i <= t.Total; i++ {
	//	q = append(q, target_broker.Task{
	//		TID:    tt.ID,
	//		Cost:   t.Cost,
	//		Title:  t.Title,
	//		Status: t.Status,
	//	})
	//}
	//
	//s.lineBroker <- q
	return nil
}
