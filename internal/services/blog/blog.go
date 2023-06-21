package blog

import (
	"core/internal/models"
	"core/internal/repositories"
)

type Service struct {
	repo *repositories.Repositories
}

func (s Service) AddComment(c models.Comment) {
	s.repo.Blog.AddComment(c)
}

func (s Service) GetBlog() []models.BlogService {
	var bss []models.BlogService
	b := s.repo.Blog.GetRecords()

	for _, v := range b {
		var bs models.BlogService
		bs.ID = v.ID
		bs.Subject = v.Subject
		bs.Text = v.Text
		bs.UID = v.UID
		bs.Comments = s.repo.Blog.GetComments(bs.ID)

		bss = append(bss, bs)
	}

	return bss
}

func NewBlogService(repo *repositories.Repositories) *Service {
	return &Service{
		repo: repo,
	}
}
