package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Login string `json:"login"`

	MainImage  string `json:"main_image"`
	SmallImage string `json:"small_image"`

	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	NumberPhone int64  `json:"number_phone"`
	Password    string `json:"password"`

	Execute          bool `json:"execute"`
	PostRegistration bool `json:"post_registration"`

	Admin bool `json:"admin"`

	Token string `json:"token"`
}

type Feed struct {
	gorm.Model

	UID    uint   `json:"uid"`
	Title  string `json:"title"`
	Link   string `json:"link"`
	Icon   string `json:"icon"`
	Status string `json:"status"`
	Cost   int64  `json:"cost"`
}

type FileStorage struct {
	gorm.Model
	Key  string `json:"key" gorm:"index"`
	Path string `json:"path"`
	Ext  string `json:"ext"`
	Type string `json:"type"`
}

type Settings struct {
	Snow bool `json:"snow"`
}
