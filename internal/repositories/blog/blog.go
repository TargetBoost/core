package blog

import (
	"core/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r Repository) GetRecords() models.Blog {
	var q models.Blog
	r.db.Table("blogs").Find(&q)
	return q
}

func (r Repository) CreateEntry(e models.CreateBlog) {
	r.db.Table("blogs").Create(e)
}

func (r Repository) UpdateEntry(e models.UpdateBlog, id uint) {
	r.db.Table("blogs").Updates(e).Where("id = ?", id)
}

func NewBlogRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
