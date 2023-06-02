package queue

import (
	"context"
	"core/internal/models"
	"core/internal/repositories"
	"core/internal/tg/bot"
	"fmt"
	"github.com/ivahaev/go-logger"
	"time"
)

const timeChange = 6000

type Queue struct {
	Line        chan []Task
	LineAppoint chan Task
	bot         *bot.Bot
	repo        *repositories.Repositories
	ctx         context.Context
}

func New(ctx context.Context, r *repositories.Repositories, bot *bot.Bot) Queue {
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
			que := q.repo.Feed.GetTaskDISTINCT()

			for _, v := range que {
				tasksUser := q.repo.Feed.GetTaskForUserUID(uint(t.UID), v.TID)
				if len(tasksUser) > 0 {
					continue
				}
				v.UID = t.UID

				v.UpdatedAt = time.Now()
				v.Status = 1
				q.repo.Feed.UpdateTask(v)
			}
		case <-q.ctx.Done():
			return

		case <-time.Tick(time.Second * 10):
			logger.Info("Check GetTaskDISTINCTIsWork()")
			que := q.repo.Feed.GetTaskDISTINCTIsWork()
			for _, v := range que {
				logger.Info(v.ID, "changed")

				if v.UpdatedAt.After(time.Now().Add(6 * time.Minute)) {
					logger.Info("change")

					v.UID = 0
					v.UpdatedAt = time.Now()
					v.Status = 0
					q.repo.Feed.UpdateTask(v)
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
				logger.Info(fmt.Sprintf(`User %v check`, v.ID))
				if v.UpdatedAt.Before(time.Now().Add((24 * 14) * time.Hour)) {
					us := q.repo.User.GetUserByID(v.ID)
					if !us.Block {
						time.Sleep(6 * time.Second)
						members, err := q.bot.CheckMembers(v.CIDChannels, v.CIDUsers)
						if err != nil {
							logger.Error(err)
						}
						if !members {
							logger.Info(fmt.Sprintf(`User %v banned`, v.ID))
							var u models.User
							u.ID = uint(v.ID)
							u.Block = true
							u.Cause = "Вы отписались от каналов раньше чем указано в правилах"

							var mm bot.Message
							mm.CID = v.CIDUsers
							mm.Type = 120

							q.bot.TrackMessages <- mm

							q.repo.User.UpdateUser(u)

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
				q.repo.Feed.CreateTask(dq)
			}
		case <-q.ctx.Done():
			return
		}
	}
}
