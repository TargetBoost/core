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
			if update.MyChatMember != nil { // If we got a message
				log.Print(update.MyChatMember.Chat.ID)
				//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				//
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				//msg.ReplyToMessageID = update.Message.MessageID
				//
				//b.API.Send(msg)
			}
		case <-b.ctx.Done():
			return
		}
	}
}
