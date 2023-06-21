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

	UID      uint             `json:"uid"`
	Text     string           `json:"text"`
	Subject  string           `json:"subject"`
	Views    int64            `json:"views"`
	Comments []CommentService `json:"comments"`
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

type CommentService struct {
	ParentID uint `json:"parent_id"`

	MainImage string `json:"main_image"`
	Login     string `json:"login"`

	CID  uint   `json:"cid"`
	UID  uint   `json:"uid"`
	Text string `json:"text"`
}
