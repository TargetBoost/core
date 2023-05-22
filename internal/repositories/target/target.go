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

func (r *Repository) GetTargetsToAdmin() []models.Target {
	var t []models.Target
	r.db.Table("targets").Find(&t)

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
	r.db.Table("targets").Where("id = ?").UpdateColumns(&target)
}

func (r *Repository) CreateTask(queue *models.Queue) {
	r.db.Table("queues").Create(&queue)
}

func (r *Repository) GetTask() []models.Queue {
	var q []models.Queue
	r.db.Table("queues").Where("uid = 0").Limit(10).Find(&q)
	return q
}
