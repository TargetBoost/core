package models

import "gorm.io/gorm"

type Blog struct {
	gorm.Model

	UID     uint
	Text    string
	Subject string
	Views   int64
}
