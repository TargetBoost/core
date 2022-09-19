package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Login string `json:"login"`

	MainImage  string `json:"main_image"`
	SmallImage string `json:"small_image"`

	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`

	Execute bool `json:"execute"`

	Admin bool `json:"admin"`

	Token string `json:"token"`
}

type Feed struct {
	gorm.Model

	UID uint `json:"uid"`

	Title string `json:"title"`

	MainImage  string `json:"main_image"`
	SmallImage string `json:"small_image"`

	Value string `json:"value"`
}

type FileStorage struct {
	gorm.Model
	Key  string `json:"key" gorm:"index"`
	Path string `json:"path"`
	Ext  string `json:"ext"`
	Type string `json:"type"`
}
