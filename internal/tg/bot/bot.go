package bot

import (
	"context"
	"core/internal/services"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
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
			if update.MyChatMember != nil {
				log.Print(update.MyChatMember.Chat.ID)
				b.services.Storage.SetChatMembers(update.MyChatMember.Chat.ID, update.MyChatMember.Chat.Title)
			}
		case <-b.ctx.Done():
			return
		}
	}
}
