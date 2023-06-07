package storage

import (
	"core/internal/models"
	"core/internal/repositories"
	"core/internal/transport/vk/api"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ivahaev/go-logger"
	"io"
	"net/http"
)

type Service struct {
	repo *repositories.Repositories
}

func (s *Service) GetFileByKey(key string) *models.FileStorage {
	return s.repo.Storage.GetFileByKey(key)
}

func (s *Service) SetChatMembers(cid, count int64, title, userName, photoLink, bio string) {
	s.repo.Storage.SetChatMembers(cid, count, title, userName, photoLink, bio)
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

	us := s.repo.Account.GetUserByToken(token)

	if us.ID == 0 {
		return errors.New("user not found")
	}

	var u models.User
	u.Token = token
	u.VKAccessToken = t.AccessToken
	u.ID = us.ID

	vk := api.New(t.AccessToken)
	vkUser, err := vk.UsersGet(nil)
	if err != nil {
		logger.Error(err)
		return err
	}

	if vkUser != nil {
		u.VKUserID = vkUser[0].ID
		u.VKUserFirstName = vkUser[0].FirstName
		u.VKUserLastName = vkUser[0].LastName
	}

	s.repo.Account.UpdateUser(u)

	logger.Debug(vkUser)

	return nil
}

func NewStorageService(repo *repositories.Repositories) *Service {
	return &Service{
		repo: repo,
	}
}
