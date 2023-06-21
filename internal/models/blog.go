package models

import "gorm.io/gorm"

type Blog struct {
	gorm.Model

	UID     uint   `json:"uid"`
	Text    string `json:"text"`
	Subject string `json:"subject"`
	Views   int64  `json:"views"`
}

type BlogService struct {
	gorm.Model

	UID      uint      `json:"uid"`
	Text     string    `json:"text"`
	Subject  string    `json:"subject"`
	Views    int64     `json:"views"`
	Comments []Comment `json:"comments"`
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

type Comment struct {
	gorm.Model
	ParentID uint `json:"parent_id"`

	CID  uint   `json:"cid"`
	UID  uint   `json:"uid"`
	Text string `json:"text"`
}