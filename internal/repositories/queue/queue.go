package queue

import (
	"core/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetTargetsToExecutor(uid int64) []models.Queue {
	var t []models.Queue
	r.db.Table("queues").Where("uid = ?", uid).Find(&t)
	return t
}

func (r *Repository) GetTaskByID(id int64) models.Queue {
	var t models.Queue
	r.db.Table("queues").Where("id = ?", id).Find(&t)
	return t
}

func (r *Repository) CreateTask(queue *models.Queue) {
	r.db.Table("queues").Create(&queue)
}

func (r *Repository) UpdateTaskStatus(q models.Queue) {
	r.db.Table("queues").UpdateColumns(&q)
}

func (r *Repository) UpdateTask(q models.Queue) {
	var qq models.Queue
	r.db.Table("queues").Where("id = ?", q.ID).Find(&qq)
	qq.UpdatedAt = q.UpdatedAt
	qq.Status = q.Status
	qq.UID = q.UID

	r.db.Debug().Save(qq)
}

func (r *Repository) GetUniqueTask() []models.Queue {
	var q []models.Queue
	r.db.Table("queues").Select("DISTINCT ON (t_id) t_id, id").Where("uid = 0").Limit(10).Order("t_id").Find(&q)
	return q
}

func (r *Repository) GetTaskDISTINCTInWork() []models.Queue {
	var q []models.Queue
	r.db.Table("queues").Select("t_id, id, uid, updated_at").Where("uid != 0 and status = 1").Find(&q)
	return q
}

func (r *Repository) GetTaskForUserUID(uid uint, tid uint) []models.Queue {
	var q []models.Queue
	r.db.Table("queues").Where("uid = ? and t_id = ?", uid, tid).Order("t_id").Find(&q)
	return q
}

func (r *Repository) GetTaskDISTINCTIsWorkForUser(uid int64) []models.QueueToExecutors {
	var q []models.QueueToExecutors
	r.db.Table(
		"queues").Select(
		"DISTINCT ON (queues.t_id) queues.t_id, queues.status, queues.id, t.title, t.link, t.icon, t.cost",
	).Joins("inner join targets t on queues.t_id = t.id").Joins("inner join chat_members_chanels on chat_members_chanels.user_name = replace(t.link, 'https://t.me/', '') ").Where(
		"queues.uid = ? and t.status = 1", uid).Order("queues.t_id").Find(&q)

	return q
}

func (r *Repository) GetChatMembersByUserName(userName string) models.ChatMembersChanel {
	var q models.ChatMembersChanel
	r.db.Table("chat_members_chanels").Where("user_name = ?", userName).Find(&q)
	return q
}

func NewQueueRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
