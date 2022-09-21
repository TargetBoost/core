package user

import (
	"core/internal/models"
	"core/internal/repositories/user"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
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
	userService.NumberPhone = v.NumberPhone
	userService.Execute = v.Execute

	return userService
}

func (s *Service) CreateUser(user models.CreateUser) (string, error) {
	token := createToken(user.Login, user.Password, time.Now())

	user.Token = token

	if user.NumberPhone == 0 || len(user.Login) == 0 {
		return "", errors.New("bad request")
	}

	if err := s.userRepository.CreateUser(&user); err != nil {
		return "", err
	}
	return token, nil
}

func createToken(login, password string, time time.Time) string {
	hash := sha256.New()
	hash.Write([]byte(login + password + time.String()))
	sha := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	return sha
}

func (s *Service) AuthUser(user models.AuthUser) (*models.User, error) {
	u := s.userRepository.GetUserByPhoneNumberAndPassword(user.NumberPhone, user.Password)

	if u.ID == 0 {
		return nil, errors.New("error auth")
	}

	token := createToken(strconv.FormatInt(user.NumberPhone, 10), user.Password, time.Now())
	u.Token = token
	s.userRepository.UpdateUser(u)

	return &u, nil
}
