package user

import (
	"core/internal/models"
	"core/internal/queue"
	"core/internal/repositories/user"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/ivahaev/go-logger"
	"strconv"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Service struct {
	userRepository *user.Repository
	lineAppoint    chan queue.Task
}

func NewUserService(userRepository *user.Repository, lineAppoint chan queue.Task) *Service {
	return &Service{
		userRepository: userRepository,
		lineAppoint:    lineAppoint,
	}
}

func (s *Service) UpdateUserBalance(id int64, cost float64) {
	u := s.userRepository.GetUserByID(id)
	logger.Info(u.Balance)
	u.Balance = u.Balance + cost - 0.50
	logger.Info(u.Balance, cost)

	s.userRepository.UpdateUser(u)
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
		userService.Balance = strconv.Itoa(int(v.Balance))

		usersService = append(usersService, userService)
	}

	return usersService
}

func (s *Service) GetUserByID(id int64) models.UserService {
	var userService models.UserService
	v := s.userRepository.GetUserByID(id)

	task := queue.Task{UID: int64(v.ID)}
	s.lineAppoint <- task

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
	userService.Admin = v.Admin
	userService.Balance = strconv.Itoa(int(v.Balance))
	userService.Block = v.Block
	userService.Cause = v.Cause
	userService.Tg = v.Tg

	return userService
}

func (s *Service) CreateUser(user models.CreateUser) (*models.User, error) {
	token := createToken(user.Login, user.Password, time.Now())

	user.Token = token

	if user.NumberPhone == 0 || len(user.Login) == 0 {
		return nil, errors.New("bad request")
	}
	if err := s.userRepository.CreateUser(&user); err != nil {
		return nil, err
	}

	u := s.userRepository.GetUserByPhoneNumberAndPassword(user.NumberPhone, user.Password)

	return &u, nil
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
