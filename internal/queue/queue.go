package queue

import (
	"context"
	"core/internal/models"
	"core/internal/repositories"
)

type Queue struct {
	Line chan []Task
	repo *repositories.Repositories
	ctx  context.Context
}

func New(ctx context.Context, r *repositories.Repositories) Queue {
	q := Queue{
		Line: make(chan []Task, 50),
		repo: r,
		ctx:  ctx,
	}

	return q
}

func (q Queue) Broker() {
	for {
		select {
		case t := <-q.Line:
			for _, v := range t {
				dq := &models.Queue{
					TID:    v.ID,
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
