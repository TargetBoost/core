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

func (r *Repository) GetUserByLogin(login string) bool {
	var u models.User
	r.db.Table("users").Where("login = ? AND deleted_at is null", login).Find(&u)

	if u.ID != 0 {
		return true
	}
	return false
}

func (r *Repository) CreateUser(user *models.CreateUser) error {
	var u models.User

	u.Login = user.Login
	u.Password = user.Password
	u.Token = user.Token
	u.NumberPhone = user.NumberPhone
	u.Execute = user.Execute

	if r.GetUserByLogin(u.Login) {
		return errors.New("user exists")
	}
	r.db.Table("users").Create(&u)

	return nil
}
