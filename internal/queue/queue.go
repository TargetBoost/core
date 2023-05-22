package queue

import (
	"context"
	"core/internal/models"
	"core/internal/repositories"
	"time"
)

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
				v.UID = t.UID
				v.UpdatedAt = time.Now()
				q.repo.Feed.UpdateTask(v)
			}
		case <-q.ctx.Done():
			return

		default:
			que := q.repo.Feed.GetTaskDISTINCTIsWork()

			for _, v := range que {
				if v.UpdatedAt.After(time.Now()) {
					v.UID = 0
					v.UpdatedAt = time.Now()
					q.repo.Feed.UpdateTask(v)
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
