package blog

import (
	"core/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r Repository) GetCommentsByParent(id uint) models.CommentParent {
	var q models.CommentParent
	r.db.Table("comments").Where("parent_id = ?", id).Find(&q)
	return q
}

func (r Repository) AddComment(c models.Comment) {
	r.db.Table("comments").Create(&c)
}

func (r Repository) GetComments(id uint) []models.Comment {
	var q []models.Comment
	r.db.Table("comments").Where("c_id = ?", id).Limit(20).Find(&q)
	return q
}

func (r Repository) GetRecords() []models.Blog {
	var q []models.Blog
	r.db.Table("blogs").Find(&q)
	return q
}

func (r Repository) CreateEntry(e models.CreateBlog) {
	r.db.Table("blogs").Create(&e)
}

func (r Repository) UpdateEntry(e models.UpdateBlog, id uint) {
	r.db.Table("blogs").Updates(e).Where("id = ?", id)
}

func NewBlogRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
