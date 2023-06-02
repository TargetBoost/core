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

				if time.Now().Unix()-v.UpdatedAt.Unix()*-1 > timeChange {
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

func (q Queue) DefenderBlocking() {
	//	select cmc.c_id, u.tg, u.id, t.link, cc.c_id from users u inner join chat_members_chanels cmc on REPLACE(u.tg, '@', '') = cmc.user_name inner join queues q on u.id = q.uid inner join targets t on q.t_id = t.id right join chat_members_chanels cc on cc.user_name = replace(t.link, 'https://t.me/', '') where q.status = 3 order by u.id

	for {
		select {
		case <-time.Tick(1 * time.Second):
			d := q.repo.Storage.GetStatisticTargetsOnExecutesIsTrue()
			for _, v := range d {
				logger.Info(v.CIDChannels, int64(v.CIDUsers)
				members, err := q.bot.CheckMembers(v.CIDChannels, int64(v.CIDUsers))
				if err != nil {
					logger.Error(err)
				}

				logger.Info(fmt.Sprintf(`User %v check`, v.ID))
				if v.UpdatedAt.After(time.Now().Add((24 * 14) * time.Hour)) {
					logger.Info(fmt.Sprintf(`User %v banned`, v.ID))
					if !members {
						us := q.repo.User.GetUserByID(v.ID)

						if !us.Block {
							var u models.User
							u.ID = uint(v.ID)
							u.Block = true
							u.Cause = "Вы отписались от каналов раньше чем указано в правилах"

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
