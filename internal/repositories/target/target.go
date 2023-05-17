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

func (r *Repository) GetTargetsToExecutor() []models.TargetToExecutors {
	var t []models.TargetToExecutors
	r.db.Table("targets t").Joins("inner join target_to_executors tte on targets.id != tte.t_id").Find(&t)

	return t
}

func (r *Repository) CreateTarget(target *models.Target) {
	r.db.Table("targets").Create(&target)
}
