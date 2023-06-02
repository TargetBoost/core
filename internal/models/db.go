package models

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type UserSettings struct {
	UID string `json:"uid"`
}

type Queue struct {
	gorm.Model

	TID    uint    `json:"tid"`
	UID    int64   `json:"uid"`
	Cost   float64 `json:"cost"`
	Title  string  `json:"title"`
	Status int64   `json:"status"`
}

type QueueToExecutors struct {
	gorm.Model

	TID    uint    `json:"tid"`
	UID    int64   `json:"uid"`
	Cost   float64 `json:"cost"`
	Title  string  `json:"title"`
	Status int64   `json:"status"`
	Icon   string  `json:"icon"` // иконка задания
	Total  float64 `json:"total"`
	Link   string  `json:"link"` // ссылка на задание
}

type ChatMembersChanel struct {
	CID       int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
	Title     string
	UserName  string
}

type QueueToService struct {
	ID     uint    `json:"id"`
	TID    uint    `json:"tid"`
	UID    int64   `json:"uid"`
	Cost   float64 `json:"cost"`
	Title  string  `json:"title"`
	Status int64   `json:"status"`
	Icon   string  `json:"icon"` // иконка задания
	Total  float64 `json:"total"`
	Link   string  `json:"link"` // ссылка на задание
}

type Transaction struct {
	gorm.Model

	BuildID string
	UID     uint
	Amount  float64
	Status  string
}

type User struct {
	gorm.Model

	Login string `json:"login"`

	Tg string `json:"tg" gorm:"index:idx_name,unique"`

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
	Block bool `json:"block"`

	Token   string  `json:"token"`
	Balance float64 `json:"balance"`
	Cause   string  `json:"cause"`
}

type Target struct {
	gorm.Model

	// основные данные
	UID        uint    `json:"uid"`    // кто создал задачу
	Title      string  `json:"title"`  // заголовок
	Link       string  `json:"link"`   // ссылка на задание
	Icon       string  `json:"icon"`   // иконка задания
	Status     int64   `json:"status"` // открыта/закрыта
	Count      int64   `json:"count"`  // количетсво заданий
	Total      float64 `json:"total"`
	Cost       float64 `json:"cost"` // цена одного задания
	TotalPrice float64 `json:"total_price"`

	// гео данные
	Country string `json:"country"`  // список стран исполнителей
	City    string `json:"city"`     // список городов исполнителей
	OldFrom int64  `json:"old_from"` // возраст исполнителя от
	OldTo   int64  `json:"old_to"`   // возраст исполнителя до
	Gender  string `json:"gender"`   // половой признак
	Type    string `json:"type"`
	Cause   string `json:"cause"`
}

type TargetToAdmin struct {
	gorm.Model

	// основные данные
	UID        uint    `json:"uid"`    // кто создал задачу
	Title      string  `json:"title"`  // заголовок
	Link       string  `json:"link"`   // ссылка на задание
	Icon       string  `json:"icon"`   // иконка задания
	Status     int64   `json:"status"` // открыта/закрыта
	Count      int64   `json:"count"`  // количетсво заданий
	Total      float64 `json:"total"`
	Cost       float64 `json:"cost"` // цена одного задания
	TotalPrice float64 `json:"total_price"`
	Login      string  `json:"login"`

	// гео данные
	Country string `json:"country"`  // список стран исполнителей
	City    string `json:"city"`     // список городов исполнителей
	OldFrom int64  `json:"old_from"` // возраст исполнителя от
	OldTo   int64  `json:"old_to"`   // возраст исполнителя до
	Gender  string `json:"gender"`   // половой признак
	Type    string `json:"type"`
	Cause   string `json:"cause"`
}

type TargetToExecutors struct {
	gorm.Model

	Status string `json:"status"`
	UID    uint   `json:"uid"`
	TID    uint   `json:"tid"`
}

type FileStorage struct {
	gorm.Model
	Key  string `json:"key" gorm:"index"`
	Path string `json:"path"`
	Ext  string `json:"ext"`
	Type string `json:"type"`
}

type Settings struct {
	ID   int64 `json:"id"`
	Snow bool  `json:"snow"`
	Rain bool  `json:"rain"`
}

type TaskCash struct {
	gorm.Model
	UID           uint
	TransactionID string
	Total         float64
	Number        string
	Status        int64
}

type UsersAndChannelsData struct {
	CIDUsers    int64
	Tg          string
	UID         uint
	Link        string
	CIDChannels int64
}

type StatisticTargetsOnExecutesIsTrue struct {
	CIDUsers    int64     `json:"cid_users" gorm:"column:cid_users"`
	CIDChannels int64     `json:"cid_channels" gorm:"column:cid_channels"`
	Tg          string    `json:"tg"`
	ID          int64     `json:"id"`
	Link        string    `json:"link"`
	UpdatedAt   time.Time `json:"updated_at"`
}
