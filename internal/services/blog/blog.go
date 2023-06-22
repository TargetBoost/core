package blog

import (
	"core/internal/models"
	"core/internal/repositories"
	"core/internal/target_broker"
	"core/internal/transport/tg/bot"
	"strings"
)

type Service struct {
	repo *repositories.Repositories

	lineAppoint   chan target_broker.Task
	trackMessages chan bot.Message
}

func (s Service) AddComment(c models.Comment, token string) {
	user := s.repo.Account.GetUserByToken(token)

	if c.ParentID != 0 {
		parentComment := s.repo.Blog.GetCommentsByParent(c.ParentID)
		userByID := s.repo.Account.GetUserByID(int64(parentComment.UID))

		st := strings.ToLower(strings.Split(userByID.Tg, "@")[len(strings.Split(userByID.Tg, "@"))-1])
		chat := s.repo.Queue.GetChatMembersByUserName(st)

		s.trackMessages <- bot.Message{
			CID:  chat.CID,
			Type: 1337,
		}
	}

	c.UID = user.ID
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
			comment.ID = vc.ID
			comment.MainImage = chat.PhotoLink
			comment.Login = user.Tg
			comment.Text = vc.Text

			if vc.ParentID != 0 {
				parentComment := s.repo.Blog.GetCommentsByParent(vc.ParentID)
				userParent := s.repo.Account.GetUserByID(int64(parentComment.UID))

				st := strings.ToLower(strings.Split(userParent.Tg, "@")[len(strings.Split(userParent.Tg, "@"))-1])
				chat := s.repo.Queue.GetChatMembersByUserName(st)

				comment.Parent = &models.CommentParent{
					CID:       parentComment.CID,
					UID:       parentComment.UID,
					Text:      parentComment.Text,
					Login:     userParent.Tg,
					MainImage: chat.PhotoLink,
				}
			}

			comments = append(comments, comment)

		}

		bs.Comments = comments

		bss = append(bss, bs)

	}

	return bss
}

func NewBlogService(repo *repositories.Repositories, lineAppoint chan target_broker.Task, trackMessages chan bot.Message) *Service {
	return &Service{
		repo: repo,

		lineAppoint:   lineAppoint,
		trackMessages: trackMessages,
	}
}
