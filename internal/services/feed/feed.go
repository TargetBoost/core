package feed

import (
	"core/internal/repositories/feed"
)

type Service struct {
	feedRepository *feed.Repository
}

func NewFeedService(feedRepository *feed.Repository) *Service {
	return &Service{
		feedRepository: feedRepository,
	}
}
