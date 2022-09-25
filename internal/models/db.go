package models

import (
	"gorm.io/gorm"
)

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

	Token   string  `json:"token"`
	Balance float64 `json:"balance"`
}

type Target struct {
	gorm.Model

	// основные данные
	UID    uint   `json:"uid"`    // кто создал задачу
	Title  string `json:"title"`  // заголовок
	Link   string `json:"link"`   // ссылка на задание
	Icon   string `json:"icon"`   // иконка задания
	Status string `json:"status"` // открыта/закрыта
	Count  int64  `json:"count"`  // количетсво заданий
	Cost   int64  `json:"cost"`   // цена одного задания

	// гео данные
	Country string `json:"country"`  // список стран исполнителей
	City    string `json:"city"`     // список городов исполнителей
	OldFrom int64  `json:"old_from"` // возраст исполнителя от
	OldTo   int64  `json:"old_to"`   // возраст исполнителя до
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
	Rain bool `json:"rain"`
}
