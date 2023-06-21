package models

import "gorm.io/gorm"

type Blog struct {
	gorm.Model

	UID     uint
	Text    string
	Subject string
	Views   int64
}

type CreateBlog struct {
	UID     uint
	Text    string
	Subject string
	Views   int64
}

type UpdateBlog struct {
	Text    string
	Subject string
	Views   int64
}
