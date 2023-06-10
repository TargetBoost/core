package target_broker

import (
	"context"
	"core/internal/models"
	"core/internal/repositories"
	bot2 "core/internal/transport/tg/bot"
	"fmt"
	"github.com/ivahaev/go-logger"
	"time"
)

const timeChange = 6000

type Queue struct {
	Line        chan []Task
	LineAppoint chan Task
	bot         *bot2.Bot
	repo        *repositories.Repositories
	ctx         context.Context
}

func New(ctx context.Context, r *repositories.Repositories, bot *bot2.Bot) Queue {
	q := Queue{
		Line:        make(chan []Task, 50),
		LineAppoint: make(chan Task, 50),
		bot:         bot,
		repo:        r,
		ctx:         ctx,
	}

	return q
}

func (q Queue) AppointTask() {
	for {
		select {
		case t := <-q.LineAppoint:
			que := q.repo.Queue.GetUniqueTask()

			for _, v := range que {
				tasksUser := q.repo.Queue.GetTaskForUserUID(uint(t.UID), v.TID)
				if len(tasksUser) > 0 {
					continue
				}
				v.UID = t.UID

				v.UpdatedAt = time.Now()
				v.Status = 1
				q.repo.Queue.UpdateTask(v)
			}
		case <-q.ctx.Done():
			return
		case <-time.Tick(time.Second * 10):
			logger.Info("Check GetTaskDISTINCTInWork()")
			que := q.repo.Queue.GetTaskDISTINCTInWork()
			for _, v := range que {
				if v.UpdatedAt.After(time.Now().Add(6 * time.Minute)) {
					logger.Info(v.ID, "changed")

					v.UID = 0
					v.UpdatedAt = time.Now()
					v.Status = 0
					q.repo.Queue.UpdateTask(v)
				}
			}
		}
	}

}

// AntiFraud - check users if unsubscribe channels BANNED
func (q Queue) AntiFraud() {
	for {
		select {
		case <-time.Tick(5 * time.Hour):
			d := q.repo.Storage.GetStatisticTargetsOnExecutesIsTrue()
			for _, v := range d {
				logger.Info(fmt.Sprintf(`Account %v check`, v.ID))
				if v.UpdatedAt.Before(time.Now().Add((24 * 14) * time.Hour)) {
					us := q.repo.Account.GetUserByID(v.ID)
					if !us.Block {
						time.Sleep(6 * time.Second)
						members, err := q.bot.CheckMembers(v.CIDChannels, v.CIDUsers)
						if err != nil {
							logger.Error(err)
						}
						if !members {
							logger.Info(fmt.Sprintf(`Account %v banned`, v.ID))
							var u models.User
							u.ID = uint(v.ID)
							u.Block = true
							u.Cause = "Вы отписались от каналов раньше чем указано в правилах"

							// Send task Message to bot sender
							q.bot.TrackMessages <- bot2.Message{
								CID:  v.CIDUsers,
								Type: 120,
							}

							q.repo.Account.UpdateUser(u)
						}
					}
				}
			}

		}
	}
}

func (q Queue) Broker() {
	for {
		select {
		case t := <-q.Line:
			for _, v := range t {
				dq := &models.Queue{
					TID:    v.TID,
					UID:    v.UID,
					Cost:   v.Cost,
					Title:  v.Title,
					Status: v.Status,
				}
				q.repo.Queue.CreateTask(dq)
			}
		case <-q.ctx.Done():
			return
		}
	}
}
