package account

import (
	"core/internal/models"
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) IsAdmin(token string) bool {
	var u models.User
	r.db.Table("users").Where("token = ? and admin = true", token).Find(&u)

	if u.Admin == true {
		return true
	}
	return false
}

func (r *Repository) GetAllUsers() []models.User {
	var u []models.User
	r.db.Table("users").Where("deleted_at is null").Order("id").Find(&u)

	return u
}

func (r *Repository) GetUserByID(id int64) models.User {
	var u models.User
	r.db.Table("users").Where("id = ? AND deleted_at is null", id).Find(&u)

	return u
}

func (r *Repository) GetUserByToken(token string) models.User {
	var u models.User
	r.db.Table("users").Where("token = ? AND deleted_at is null", token).Find(&u)

	return u
}

func (r *Repository) GetUserByTgAndPassword(tg, pass string) models.User {
	var u models.User
	r.db.Table("users").Where("tg = ? AND password = ? AND deleted_at is null", tg, pass).Find(&u)

	return u
}

func (r *Repository) IsUserByTg(tg string) bool {
	var u models.User
	r.db.Table("users").Where("tg = ? AND deleted_at is null", tg).Find(&u)

	if u.ID != 0 {
		return true
	}
	return false
}

func (r *Repository) UpdateUser(user models.User) {
	r.db.Debug().Table("users").Where("id = ?", user.ID).Updates(user)
}

func (r *Repository) UpdateUserBalance(uid uint, balance float64) {
	var q models.User
	r.db.Table("users").Where("id = ?", uid).Find(&q)
	//r.db.Debug().Table("users").Where("id = ?", uid).Updates(models.Account{Balance: balance})
	q.Balance = balance
	r.db.Debug().Save(q)
}

func (r *Repository) CreateUser(user *models.CreateUser) error {
	var u models.User

	if len(user.Login) > 20 {
		u.Tg = user.Tg[:20]
	} else {
		u.Tg = user.Tg

	}
	u.Password = user.Password
	u.Token = user.Token
	u.Tg = user.Tg
	//u.NumberPhone = user.NumberPhone
	u.Execute = user.Execute

	if r.IsUserByTg(u.Tg) {
		return errors.New("user exists")
	}
	r.db.Table("users").Create(&u)

	return nil
}

func (r *Repository) CreateTaskCache(task models.TaskCash) {
	r.db.Table("task_cashes").Create(&task)
}

func (r *Repository) UpdateTaskCache(task models.TaskCash) {
	r.db.Table("task_cashes").Updates(task)
}

func (r *Repository) GetTaskCacheByUID(uid uint) []models.TaskCash {
	var q []models.TaskCash
	r.db.Table("task_cashes").Where("uid = ?", uid).Order("created_at").Find(&q)
	return q
}

func (r *Repository) GetTaskCacheByID(id uint) models.TaskCash {
	var q models.TaskCash
	r.db.Table("task_cashes").Where("id = ?", id).Find(&q)
	return q
}

func (r *Repository) GetTaskCacheToAdmin() []models.TaskCash {
	var q []models.TaskCash
	r.db.Table("task_cashes").Order("created_at").Find(&q)

	return q
}

func (r *Repository) CreateTransaction(t *models.TransactionToService) {
	r.db.Table("transactions").Create(&t)
}

func (r *Repository) UpdateTransaction(t *models.TransactionToService) {
	r.db.Table("transactions").Where("build_id = ?", t.BuildID).Updates(&t)
}

func (r *Repository) GetTransaction(build string) models.Transaction {
	var q models.Transaction
	r.db.Table("transactions").Where("build_id = ?", build).Find(&q)
	return q
}

func (r *Repository) IsAuth(token string) (uint, bool) {
	var u models.User
	r.db.Table("users").Where("token = ?", token).Find(&u)

	if u.ID != 0 {
		return u.ID, true
	}

	return 0, false
}

func NewAccountRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
