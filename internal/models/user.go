package models

import (
	"time"
)

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

	Execute bool   `json:"execute"`
	Admin   bool   `json:"admin"`
	Balance string `json:"balance"`
	Block   bool   `json:"block"`
	Cause   string `json:"cause"`
	Tg      string `json:"tg"`
	Token   string `json:"token"`
	VKToken string `json:"vk_token"`
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
	Tg       string `json:"tg"`
	Password string `json:"password"`
}

type Pay struct {
	Value string `json:"value"`
}

type TransactionToService struct {
	BuildID string
	UID     uint
	Amount  string
	Status  string
}

type TaskCashToService struct {
	ID            uint   `json:"id"`
	UID           uint   `json:"uid"`
	TransactionID string `json:"transaction_id"`
	Number        string `json:"number"`
	Total         string `json:"total"`
	Status        int64  `json:"status"`
}

type TaskCashToUser struct {
	Total  float64 `json:"total"`
	Number string  `json:"number"`
}
