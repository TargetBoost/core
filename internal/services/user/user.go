package user

import (
	"core/internal/models"
	"core/internal/queue"
	"core/internal/repositories/user"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
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
	u.Balance = u.Balance + cost - 1
	logger.Info(u.Balance, cost)

	s.userRepository.UpdateUser(&u)
}

func (s *Service) UpdateUser(uid uint, b float64) {
	var u models.User
	u.Balance = b
	u.ID = uid
	s.userRepository.UpdateUser(&u)
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

func (s *Service) GetTasksCashesUser(uid uint) []models.TaskCashToService {
	tasks := func(t []models.TaskCash, f func(t models.TaskCash) models.TaskCashToService) []models.TaskCashToService {
		result := make([]models.TaskCashToService, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.userRepository.GetTaskCacheByUID(uid), models.MapToTasksUser)

	return tasks
}

func (s *Service) GetTasksCashesAdmin() []models.TaskCashToService {
	tasks := func(t []models.TaskCash, f func(t models.TaskCash) models.TaskCashToService) []models.TaskCashToService {
		result := make([]models.TaskCashToService, 0, len(t))
		for _, value := range t {
			result = append(result, f(value))
		}

		return result
	}(s.userRepository.GetTaskCacheToAdmin(), models.MapToTasksUser)

	return tasks
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
	userService.Balance = strconv.FormatFloat(v.Balance, 'g', -1, 64)
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
	s.userRepository.UpdateUser(&u)

	return &u, nil
}

func (s *Service) CreateTaskCashes(uid int64, task models.TaskCashToUser) error {
	u := s.userRepository.GetUserByID(uid)

	id := uuid.New()

	if task.Total < 1.00 {
		return errors.New("Сумма вывода не может быть меньше 1.00")
	}

	if u.Balance < 5 {
		return errors.New("Ваш баланс меньше минимального вывода")
	}

	if u.Balance < task.Total {
		return errors.New("Сумма вывода больше баланса")
	}

	u.Balance = u.Balance - task.Total

	s.userRepository.UpdateUser(&u)

	var t models.TaskCash
	t.Status = 0
	t.UID = u.ID
	t.Number = task.Number
	t.Total = task.Total
	t.TransactionID = id.String()
	s.userRepository.CreateTaskCache(t)

	return nil
}

func (s *Service) UpdateTaskCashes(task models.TaskCashToService) {
	var q models.TaskCash

	q.ID = task.ID
	q.Status = task.Status

	s.userRepository.UpdateTaskCache(q)
}

func (s *Service) CreateTransaction(t *models.TransactionToService) {
	s.userRepository.CreateTransaction(t)
}

func (s *Service) UpdateTransaction(t *models.TransactionToService) {
	s.userRepository.UpdateTransaction(t)
}

func (s *Service) GetTransaction(build string) *models.TransactionToService {
	var trans models.TransactionToService

	t := s.userRepository.GetTransaction(build)

	trans.Amount = strconv.FormatFloat(t.Amount, 'g', -1, 64)
	trans.UID = t.UID
	trans.Status = t.Status
	trans.BuildID = t.BuildID

	return &trans
}
