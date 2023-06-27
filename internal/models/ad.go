package models

import "gorm.io/gorm"

// Campaign - is base struct for ad
type Campaign struct {
	gorm.Model

	UID         uint
	Name        string
	Description string
	Link        string
	Limit       int64
	Type        Type   `gorm:"references:ID"`
	Status      Status `gorm:"references:ID"`
}

type Status struct {
	ID    uint `gorm:"primarykey"`
	Value string
	Meta  interface{}
}

type Type struct {
	ID    uint `gorm:"primarykey"`
	Value string
	Meta  interface{}
}
