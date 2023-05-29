package queue

import (
	"context"
	"core/internal/models"
	"core/internal/repositories"
	"github.com/ivahaev/go-logger"
	"time"
)

const timeChange = 20.0

type Queue struct {
	Line        chan []Task
	LineAppoint chan Task
	repo        *repositories.Repositories
	ctx         context.Context
}

func New(ctx context.Context, r *repositories.Repositories) Queue {
	q := Queue{
		Line:        make(chan []Task, 50),
		LineAppoint: make(chan Task, 50),
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

		default:
			que := q.repo.Feed.GetTaskDISTINCTIsWork()
			for _, v := range que {
				if v.UID == 0 {
					continue
				}

				if time.Now().Sub(v.UpdatedAt).Seconds() > timeChange {
					logger.Info(time.Now().Sub(v.UpdatedAt).Seconds(), timeChange)
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
		case <-time.Tick(12 * time.Hour):
			//q.repo.Storage.
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
