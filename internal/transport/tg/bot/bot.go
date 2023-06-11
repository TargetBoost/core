package bot

import (
	"context"
	"core/internal/repositories"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivahaev/go-logger"
	"image/jpeg"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/imageorient"
)

const (
	directoryPath = `./uploads/tg_chats_photos`
	filesPath     = `./uploads/tg_chats_photos/%s`
	tgFilesPath   = `https://api.telegram.org/file/bot%s/%s`
)

type Bot struct {
	API *tgbotapi.BotAPI

	token        string
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
		token:         token,
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

// GetUpdates - get updates bot messages
func (b *Bot) GetUpdates() {
	updates := b.API.GetUpdatesChan(b.updateConfig)
	for update := range updates {
		if update.MyChatMember != nil {
			logger.Info(fmt.Sprintf("New Chat ID: %v", update.MyChatMember.Chat.ID))

			if chat, err := b.API.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: update.MyChatMember.Chat.ID}}); err != nil || chat.Photo == nil {
				if err != nil {
					logger.Error(err)
					continue
				}

				count, err := b.API.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{
					ChatConfig: tgbotapi.ChatConfig{
						ChatID: update.MyChatMember.Chat.ID,
					},
				})
				if err != nil {
					logger.Error(err)
				}

				logger.Debug(chat)

				b.repos.Storage.SetChatMembers(update.MyChatMember.Chat.ID, int64(count), update.MyChatMember.Chat.Title, strings.ToLower(update.MyChatMember.Chat.UserName), "", chat.Description)
				continue
			} else {
				fileID := chat.Photo.BigFileID
				file, err := b.API.GetFile(tgbotapi.FileConfig{
					FileID: fileID,
				})
				if err != nil {
					logger.Error(err)
				}

				//logger.Info(fmt.Sprintf(tgFilesPath, b.token, file.FilePath))
				err = downloadFile(fmt.Sprintf(filesPath, file.FileID), fmt.Sprintf(tgFilesPath, b.token, file.FilePath))
				if err != nil {
					logger.Error(err)
				}

				count, err := b.API.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{
					ChatConfig: tgbotapi.ChatConfig{
						ChatID: update.MyChatMember.Chat.ID,
					},
				})
				if err != nil {
					logger.Error(err)
				}

				logger.Debug(chat)
				b.repos.Storage.SetChatMembers(update.MyChatMember.Chat.ID, int64(count), update.MyChatMember.Chat.Title, strings.ToLower(update.MyChatMember.Chat.UserName), file.FileID, chat.Description)
			}
		}
		if update.Message != nil {
			if chat, err := b.API.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: update.Message.Chat.ID}); err != nil {
				logger.Debug(chat.Photos)
				if len(chat.Photos) > 0 {
					fileID := chat.Photos[0][0].FileID
					file, err := b.API.GetFile(tgbotapi.FileConfig{
						FileID: fileID,
					})
					if err != nil {
						logger.Error(err)
					}

					//logger.Info(fmt.Sprintf(tgFilesPath, b.token, file.FilePath))
					err = downloadFile(fmt.Sprintf(filesPath, file.FileID), fmt.Sprintf(tgFilesPath, b.token, file.FilePath))
					if err != nil {
						logger.Error(err)
					}

					b.repos.Storage.SetChatMembers(update.Message.Chat.ID, int64(0), update.Message.Chat.Title, strings.ToLower(update.Message.Chat.UserName), file.FileID, "")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, `
Добро пожаловать!
Вы добавлены в систему.
				`)
					b.API.Send(msg)
					continue
				}
			}
			logger.Info(update.Message.Chat)

			b.repos.Storage.SetChatMembers(update.Message.Chat.ID, int64(0), update.Message.Chat.Title, strings.ToLower(update.Message.Chat.UserName), "", "")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `
Добро пожаловать!
Вы добавлены в систему.
				`)
			b.API.Send(msg)
		}
	}
}

func downloadFile(filepath string, url string) error {
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err := os.Mkdir(directoryPath, 0777)
		return err
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	img, _, err := imageorient.Decode(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("imageorient.Decode failed: %v", err))
	}

	err = jpeg.Encode(out, img, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("jpeg.Encode failed: %v", err))
	}

	//// Writer the body to file
	//_, err = io.Copy(out, resp.Body)
	//if err != nil {
	//	return err
	//}

	return nil
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
