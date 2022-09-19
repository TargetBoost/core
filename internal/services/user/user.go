package user

import (
	"core/internal/models"
	"core/internal/repositories/user"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Service struct {
	userRepository *user.Repository
}

func NewUserService(userRepository *user.Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) GetAllUsers() []models.UserService {
	var usersService []models.UserService

	for _, v := range s.userRepository.GetAllUsers() {
		var userService models.UserService

		userService.ID = v.ID
		userService.CreatedAt = v.CreatedAt
		userService.Login = v.Login
		userService.FirstName = v.FirstName
		userService.LastName = v.LastName
		userService.MiddleName = v.MiddleName
		userService.MainImage = v.MainImage
		userService.SmallImage = v.SmallImage

		usersService = append(usersService, userService)
	}

	return usersService
}

func (s *Service) GetUserByID(id int64) models.UserService {
	var userService models.UserService
	v := s.userRepository.GetUserByID(id)

	userService.ID = v.ID
	userService.CreatedAt = v.CreatedAt
	userService.Login = v.Login
	userService.FirstName = v.FirstName
	userService.LastName = v.LastName
	userService.MiddleName = v.MiddleName
	userService.MainImage = v.MainImage
	userService.SmallImage = v.SmallImage

	return userService
}

func (s *Service) CreateUser(user models.CreateUser) error {
	token := createToken(user.Login, user.Password, time.Now())

	user.Token = token

	if err := s.userRepository.CreateUser(&user); err != nil {
		return err
	}
	return nil
}

func createToken(login, password string, time time.Time) string {
	hash := sha256.New()
	hash.Write([]byte(login + password + time.String()))
	sha := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	return sha
}
