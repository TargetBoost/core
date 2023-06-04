package storage

import (
	"core/internal/models"
	"core/internal/repositories/storage"
	"core/internal/repositories/user"
	"core/internal/vk/api"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ivahaev/go-logger"
	"io"
	"net/http"
)

type Service struct {
	storageRepository *storage.Repository
	userRepository    *user.Repository
}

func NewStorageService(storageRepository *storage.Repository, userRepository *user.Repository) *Service {
	return &Service{
		storageRepository: storageRepository,
		userRepository:    userRepository,
	}
}

func (s *Service) GetFileByKey(key string) *models.FileStorage {
	return s.storageRepository.GetFileByKey(key)
}

func (s *Service) SetChatMembers(cid, count int64, title, userName, photoLink, bio string) {
	s.storageRepository.SetChatMembers(cid, count, title, userName, photoLink, bio)
}

func (s *Service) CallBackVK(code, token string) error {
	httpClient := http.Client{}

	requestURL := fmt.Sprintf("https://oauth.vk.com/access_token?client_id=51666148&client_secret=vvCXlyIJ0yEIkOyAHrAV&redirect_uri=https://targetboost.ru/core/v1/callback_vk&code=%s", code)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		logger.Errorf("could not create HTTP request: %v", err)
		return err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		logger.Errorf("could not send HTTP request: %v", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error(err)
			return
		}
	}(res.Body)

	type Result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		UserId      int    `json:"user_id"`
	}

	var t Result

	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		logger.Errorf("could not parse JSON response: %v", err)
		return err
	}

	us := s.userRepository.GetUserByToken(token)

	if us.ID == 0 {
		return errors.New("user not found")
	}

	var u models.User
	u.Token = token
	u.VKAccessToken = t.AccessToken
	u.ID = us.ID

	s.userRepository.UpdateUser(u)

	vk := api.New(t.AccessToken)
	vkUser, err := vk.UsersGet(nil)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Debug(vkUser)

	return nil
}
