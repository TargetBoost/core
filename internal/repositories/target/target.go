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

	var targetResult []models.Target

	for _, v := range t {
		var sc models.SubCount

		r.db.Table("queues").Select("count(t_id)").Where("t_id = ? and status = 3", v.ID).Find(&sc)

		v.Count = sc.Count

		targetResult = append(targetResult, v)
	}

	return targetResult
}

func (r *Repository) GetTargetByID(uid uint) models.Target {
	var t models.Target
	r.db.Table("targets").Where("id = ?", uid).Find(&t)

	return t
}

func (r *Repository) GetTargetsToAdmin() []models.TargetToAdmin {
	var t []models.TargetToAdmin
	r.db.Table("targets").Select("targets.id, targets.uid, targets.title, targets.link, targets.icon, targets.status, targets.count, targets.total, targets.cost, targets.total_price, u.login").Joins("inner join users u on targets.uid = u.id").Find(&t)

	var targetResult []models.TargetToAdmin

	for _, v := range t {
		var sc models.SubCount

		r.db.Table("queues").Select("count(t_id)").Where("t_id = ? and status = 3", v.ID).Find(&sc)

		v.Count = sc.Count

		targetResult = append(targetResult, v)
	}

	return targetResult
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

func (r *Repository) GetTaskDISTINCT() []models.Queue {
	var q []models.Queue
	r.db.Table("queues").Select("DISTINCT ON (t_id) t_id, id").Where("uid = 0").Limit(10).Order("t_id").Find(&q)
	return q
}

func (r *Repository) GetTaskDISTINCTIsWork() []models.Queue {
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
	r.db.Table("queues").Select("DISTINCT ON (queues.t_id) queues.t_id, queues.status, queues.id, t.title, t.link, t.icon, t.cost").Joins("inner join targets t on queues.t_id = t.id").Where("queues.uid = ? and t.status = 1", uid).Order("queues.t_id").Find(&q)
	return q
}

func (r *Repository) GetChatMembersByUserName(userName string) models.ChatMembersChanel {
	var q models.ChatMembersChanel
	r.db.Table("chat_members_chanels").Where("user_name = ?", userName).Find(&q)
	return q
}
