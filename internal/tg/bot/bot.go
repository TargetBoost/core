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
	if err != nil {
		return nil, err
	}

	api.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return &Bot{
		API:          api,
		updateConfig: u,
		ctx:          ctx,
		services:     services,
	}, nil
}

//func (b *Bot) SenderUpdates(){
//	b.services.Storage
//}

func (b *Bot) GetUpdates() {
	updates := b.API.GetUpdatesChan(b.updateConfig)
	for update := range updates {
		//logger.Info(update.Message.Chat)
		if update.MyChatMember != nil {
			logger.Info(update.MyChatMember)
			b.services.Storage.SetChatMembers(update.MyChatMember.Chat.ID, update.MyChatMember.Chat.Title, update.MyChatMember.Chat.UserName)
		}
		if update.Message != nil {
			logger.Info(update.Message.Chat)
			b.services.Storage.SetChatMembers(update.Message.Chat.ID, update.Message.Chat.Title, update.Message.Chat.UserName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `
Добро пожаловать!
Вы добавлены в систему.
				`)
			b.API.Send(msg)
		}
	}
	//	for {
	//		select {
	//		case update := <-b.API.GetUpdatesChan(b.updateConfig):
	//			//logger.Info(update.Message.Chat)
	//			if update.MyChatMember != nil {
	//				logger.Info(update.MyChatMember)
	//				b.services.Storage.SetChatMembers(update.MyChatMember.Chat.ID, update.MyChatMember.Chat.Title, update.MyChatMember.Chat.UserName)
	//			}
	//			if update.Message != nil {
	//				logger.Info(update.Message.Chat)
	//				b.services.Storage.SetChatMembers(update.Message.Chat.ID, update.Message.Chat.Title, update.Message.Chat.UserName)
	//				msg := tgbotapi.NewMessage(update.Message.Chat.ID, `
	//Добро пожаловать!
	//Вы добавлены в систему.
	//				`)
	//				b.API.Send(msg)
	//			}
	//
	//		case <-b.ctx.Done():
	//			return
	//		}
	//	}
}

func (b *Bot) CheckMembers(cid, uid int64) (bool, error) {
	member, err := b.API.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: cid,
			UserID: uid,
		},
	})
	if err != nil {
		return false, err
	}

	if member.Status == "member" {
		return true, nil
	}
	return false, nil
}
