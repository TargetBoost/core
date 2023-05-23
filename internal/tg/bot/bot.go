package bot

import (
	"context"
	"core/internal/services"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivahaev/go-logger"
)

type Bot struct {
	API *tgbotapi.BotAPI

	services     *services.Services
	updateConfig tgbotapi.UpdateConfig
	ctx          context.Context
}

func New(ctx context.Context, token string, services *services.Services) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	api.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	if err != nil {
		return nil, err
	}

	return &Bot{
		API:          api,
		updateConfig: u,
		ctx:          ctx,
		services:     services,
	}, nil
}

func (b *Bot) GetUpdates() {
	for {
		select {
		case update := <-b.API.GetUpdatesChan(b.updateConfig):
			logger.Info(update.Message.Chat)
			if update.Message != nil {
				b.services.Storage.SetChatMembers(update.Message.Chat.ID, update.Message.Chat.Title, update.Message.Chat.UserName)
			}
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *Bot) CheckMembers(cid, uid int64) error {
	_, err := b.API.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: cid,
			UserID: uid,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
