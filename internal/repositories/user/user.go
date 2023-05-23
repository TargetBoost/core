package user

import (
	"core/internal/models"
	"github.com/kataras/iris/v12/x/errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAllUsers() []models.User {
	var u []models.User
	r.db.Table("users").Where("deleted_at is null").Find(&u)

	return u
}

func (r *Repository) GetUserByID(id int64) models.User {
	var u models.User
	r.db.Table("users").Where("id = ? AND deleted_at is null", id).Find(&u)

	return u
}

func (r *Repository) GetUserByPhoneNumberAndPassword(ph int64, pass string) models.User {
	var u models.User
	r.db.Table("users").Where("number_phone = ? AND password = ? AND deleted_at is null", ph, pass).Find(&u)

	return u
}

func (r *Repository) GetUserByLogin(login string) bool {
	var u models.User
	r.db.Table("users").Where("login = ? AND deleted_at is null", login).Find(&u)

	if u.ID != 0 {
		return true
	}
	return false
}

func (r *Repository) UpdateUser(user *models.User) {
	r.db.Table("users").Updates(user)
}

func (r *Repository) CreateUser(user *models.CreateUser) error {
	var u models.User

	if len(user.Login) > 20 {
		u.Login = user.Login[:20]
	} else {
		u.Login = user.Login

	}
	u.Password = user.Password
	u.Token = user.Token
	u.Tg = user.Tg
	u.NumberPhone = user.NumberPhone
	u.Execute = user.Execute

	if r.GetUserByLogin(u.Login) {
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
