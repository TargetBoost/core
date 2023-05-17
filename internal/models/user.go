package models

import "time"

type UserService struct {
	ID uint `json:"id"`

	CreatedAt time.Time `json:"created_at"`

	Login string `json:"login"`

	MainImage  string `json:"main_image"`
	SmallImage string `json:"small_image"`

	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`

	NumberPhone int64 `json:"number_phone"`

	Execute bool  `json:"execute"`
	Admin   bool  `json:"admin"`
	Balance int64 `json:"balance"`
}

type CreateUser struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	NumberPhone int64  `json:"number_phone"`
	Execute     bool   `json:"execute"`
	Tg          string `json:"tg"`
}

type AuthUser struct {
	NumberPhone int64  `json:"number_phone"`
	Password    string `json:"password"`
}
