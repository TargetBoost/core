package bot

import (
	"context"
	"core/internal/repositories"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivahaev/go-logger"
	"strings"
)

type Bot struct {
	API *tgbotapi.BotAPI

	repos        *repositories.Repositories
	updateConfig tgbotapi.UpdateConfig
	ctx          context.Context

	TrackMessages chan Message
}

func New(ctx context.Context, token string, repos *repositories.Repositories) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return &Bot{
		API:           api,
		updateConfig:  u,
		ctx:           ctx,
		repos:         repos,
		TrackMessages: make(chan Message, 100),
	}, nil
}

func (b *Bot) SenderUpdates() {
	for {
		select {
		case m := <-b.TrackMessages:
			switch m.Type {
			case 4:
				msg := tgbotapi.NewMessage(m.CID, fmt.Sprintf(`Ваша заявка на вывод средств отклонена (%vруб.)`, m.Count))
				_, err := b.API.Send(msg)
				if err != nil {
					logger.Error(err)
				}
			case 2:
				msg := tgbotapi.NewMessage(m.CID, fmt.Sprintf(`Деньги по Вашей заявке успешго отправлены (%vруб.)`, m.Count))
				_, err := b.API.Send(msg)
				if err != nil {
					logger.Error(err)
				}
			case 100:
				msg := tgbotapi.NewMessage(m.CID, `Для Вас появились новые задания!`)
				_, err := b.API.Send(msg)
				if err != nil {
					logger.Error(err)
				}
			case 120:
				msg := tgbotapi.NewMessage(m.CID, `Ваш профиль был заблокирован за отписку от каналов раньше 2 недель.`)
				_, err := b.API.Send(msg)
				if err != nil {
					logger.Error(err)
				}
			default:
				msg := tgbotapi.NewMessage(m.CID, fmt.Sprintf(`Заявка на вывод средств создана (%vруб.)`, m.Count))
				_, err := b.API.Send(msg)
				if err != nil {
					logger.Error(err)
				}
			}

		case <-b.ctx.Done():
			break
		}
	}

}

func (b *Bot) GetUpdates() {
	updates := b.API.GetUpdatesChan(b.updateConfig)
	for update := range updates {
		//logger.Info(update.Message.Chat)
		if update.MyChatMember != nil {
			logger.Info(update.MyChatMember)

			if photos, err := b.API.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: update.MyChatMember.Chat.ID}); err != nil || len(photos.Photos) == 0 {
				logger.Error(err)
				b.repos.Storage.SetChatMembers(update.MyChatMember.Chat.ID, update.MyChatMember.Chat.Title, strings.ToLower(update.MyChatMember.Chat.UserName), "")
				continue
			} else {
				fileID := photos.Photos[0][0].FileID
				file, err := b.API.GetFile(tgbotapi.FileConfig{
					FileID: fileID,
				})
				if err != nil {
					logger.Error(err)
				}

				b.repos.Storage.SetChatMembers(update.MyChatMember.Chat.ID, update.MyChatMember.Chat.Title, strings.ToLower(update.MyChatMember.Chat.UserName), file.FilePath)
			}
		}
		if update.Message != nil {
			logger.Info(update.Message.Chat)

			b.repos.Storage.SetChatMembers(update.Message.Chat.ID, update.Message.Chat.Title, strings.ToLower(update.Message.Chat.UserName), "")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `
Добро пожаловать!
Вы добавлены в систему.
				`)
			b.API.Send(msg)
		}
	}
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

	//logger.Info(member)

	if member.Status == "member" {
		return true, nil
	}
	return false, nil
}
