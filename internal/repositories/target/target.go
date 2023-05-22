package target

import (
	"core/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewTargetRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetTargets(uid uint) []models.Target {
	var t []models.Target
	r.db.Table("targets").Where("uid = ?", uid).Find(&t)

	return t
}

func (r *Repository) GetTargetByID(uid uint) models.Target {
	var t models.Target
	r.db.Table("targets").Where("id = ?", uid).Find(&t)

	return t
}

func (r *Repository) GetTargetsToAdmin() []models.TargetToAdmin {
	var t []models.TargetToAdmin
	r.db.Table("targets").Select("targets.id, targets.uid, targets.title, targets.link, targets.icon, targets.status, targets.count, targets.total, targets.cost, targets.total_price, u.login").Joins("inner join users u on targets.uid = u.id").Find(&t)

	return t
}

func (r *Repository) GetTargetsToExecutor(uid int64) []models.Queue {
	var t []models.Queue
	r.db.Table("queues").Where("uid = ?", uid).Find(&t)
	return t
}

func (r *Repository) CreateTarget(target *models.Target) *models.Target {
	r.db.Table("targets").Create(&target)
	return target
}

func (r *Repository) UpdateTarget(id uint, target *models.Target) {
	r.db.Table("targets").Where("id = ?", id).UpdateColumns(&target)
}

func (r *Repository) CreateTask(queue *models.Queue) {
	r.db.Table("queues").Create(&queue)
}

func (r *Repository) UpdateTask(q models.Queue) {
	r.db.Table("queues").UpdateColumns(&q).Where("uid = 0")
}

func (r *Repository) GetTaskDISTINCT() []models.Queue {
	var q []models.Queue
	r.db.Table("queues").Select("DISTINCT ON (t_id) t_id, id").Where("uid = 0").Limit(10).Order("t_id").Find(&q)
	return q
}

func (r *Repository) GetTaskDISTINCTIsWork() []models.Queue {
	var q []models.Queue
	r.db.Table("queues").Select("t_id, id").Where("uid != 0 and status != 3").Order("t_id").Find(&q)
	return q
}

func (r *Repository) GetTaskDISTINCTIsWorkForUser(uid int64) []models.QueueToExecutors {
	var q []models.QueueToExecutors
	r.db.Table("queues").Select("DISTINCT ON (queues.t_id) queues.t_id, queues.status, queues.id, t.title, t.link, t.icon, t.cost").Joins("inner join targets t on queues.t_id = t.id").Where("queues.uid = ? and t.status = 1", uid).Order("queues.t_id").Find(&q)
	return q
}
