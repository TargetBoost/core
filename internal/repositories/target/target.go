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

func (r *Repository) GetTargetsToAdmin() []models.Target {
	var t []models.Target
	r.db.Table("targets").Find(&t)

	return t
}

func (r *Repository) GetTargetsToExecutor() []models.Target {
	var t []models.Target
	r.db.Table("targets").Select("targets.uid, targets.created_at, targets.updated_at, targets.deleted_at, targets.status,  targets.count, targets.cost, targets.total, targets.link, targets.icon").Joins("inner join target_to_executors on targets.id != target_to_executors.t_id").Find(&t)

	return t
}

func (r *Repository) CreateTarget(target *models.Target) {
	r.db.Table("targets").Create(&target)
}
