package blog

import (
	"core/internal/models"
	"core/internal/repositories"
	"strings"
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
		c := s.repo.Blog.GetComments(bs.ID)

		var comments []models.CommentService
		for _, vc := range c {
			var comment models.CommentService

			user := s.repo.Account.GetUserByID(int64(vc.UID))
			st := strings.ToLower(strings.Split(user.Tg, "@")[len(strings.Split(user.Tg, "@"))-1])
			chat := s.repo.Queue.GetChatMembersByUserName(st)
			comment.MainImage = chat.PhotoLink
			comment.Login = user.Tg
			comment.Text = vc.Text

			comments = append(comments, comment)
		}

		bs.Comments = comments

		bss = append(bss, bs)
	}

	return bss
}

func NewBlogService(repo *repositories.Repositories) *Service {
	return &Service{
		repo: repo,
	}
}
