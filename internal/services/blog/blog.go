package blog

import (
	"core/internal/models"
	"core/internal/repositories"
)

type Service struct {
	repo *repositories.Repositories
}

func (s Service) GetBlog() []models.Blog {
	return s.repo.Blog.GetRecords()
}

func NewBlogService(repo *repositories.Repositories) *Service {
	return &Service{
		repo: repo,
	}
}
