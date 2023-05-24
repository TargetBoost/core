package handler

import (
	"bytes"
	"core/internal/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/ivahaev/go-logger"
	"github.com/kataras/iris/v12"
	"io"
	"net/http"
)

// GetAllUsers all users returned
func (h *Handler) GetAllUsers(ctx iris.Context) {
	users := h.Service.User.GetAllUsers()

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": users,
	})
}

// CreateUser create user
func (h *Handler) CreateUser(ctx iris.Context) {
	var u models.CreateUser
	err := ctx.ReadJSON(&u)
	if err != nil {
		ctx.StatusCode(400)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "bad data insertion",
			},
			"data": nil,
		})
		return
	}

	logger.Debug(u)

	user, err := h.Service.User.CreateUser(u)
	if err != nil {
		ctx.StatusCode(200)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
		})
		return
	}

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": iris.Map{
			"token": user.Token,
			"id":    user.ID,
		},
	})
}

// CreateTaskCashes create
func (h *Handler) CreateTaskCashes(ctx iris.Context) {
	rawToken := ctx.GetHeader("Authorization")
	var task models.TaskCashToUser
	err := ctx.ReadJSON(&task)
	if err != nil {
		ctx.StatusCode(400)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "bad data insertion",
			},
			"data": nil,
		})
		return
	}

	u, err := h.CheckAuth(rawToken)
	if err != nil {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
			"data": nil,
		})
		return
	}

	err = h.Service.User.CreateTaskCashes(int64(u.ID), task)
	if err != nil {
		ctx.StatusCode(200)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
		})
		return
	}

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": iris.Map{},
	})
}

// UpdateTaskCashes create
func (h *Handler) UpdateTaskCashes(ctx iris.Context) {
	rawToken := ctx.GetHeader("Authorization")
	var task models.TaskCashToService
	err := ctx.ReadJSON(&task)
	if err != nil {
		ctx.StatusCode(400)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "bad data insertion",
			},
			"data": nil,
		})
		return
	}

	u, err := h.CheckAuth(rawToken)
	if err != nil {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
			"data": nil,
		})
		return
	}

	if !u.Admin {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "your dont have permission",
			},
			"data": nil,
		})
		return
	}

	h.Service.User.UpdateTaskCashes(task)
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": iris.Map{},
	})
}

func (h *Handler) GetTaskCashes(ctx iris.Context) {
	rawToken := ctx.GetHeader("Authorization")
	u, err := h.CheckAuth(rawToken)
	if err != nil {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
			"data": nil,
		})
		return
	}

	tasks := h.Service.User.GetTasksCashesUser(u.ID)
	if err != nil {
		ctx.StatusCode(200)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
		})
		return
	}

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": tasks,
	})
}

func (h *Handler) GetTaskCashesAdmin(ctx iris.Context) {
	rawToken := ctx.GetHeader("Authorization")
	u, err := h.CheckAuth(rawToken)
	if err != nil {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
			"data": nil,
		})
		return
	}

	if !u.Admin {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "your dont have permission",
			},
			"data": nil,
		})
		return
	}

	tasks := h.Service.User.GetTasksCashesAdmin()
	if err != nil {
		ctx.StatusCode(200)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
		})
		return
	}

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": tasks,
	})
}

// Pay
func (h *Handler) Pay(ctx iris.Context) {
	rawToken := ctx.GetHeader("Authorization")
	var pay models.Pay
	err := ctx.ReadJSON(&pay)
	user, err := h.CheckAuth(rawToken)
	if err != nil {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
			"data": nil,
		})
		return
	}

	if user.ID == 0 {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "User not exist",
			},
			"data": nil,
		})
		return
	}

	id := uuid.New()

	type amount struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	}

	type Body struct {
		Amount amount `json:"amount"`
	}

	var body Body

	body.Amount = amount{
		Currency: "RUB",
		Value:    pay.Value,
	}

	jsonBody, _ := json.Marshal(body)
	bodyReader := bytes.NewReader(jsonBody)

	httpClient := http.Client{}

	reqURL := fmt.Sprintf("https://api.qiwi.com/partner/bill/v1/bills/%s", id.String())
	req, err := http.NewRequest(http.MethodGet, reqURL, bodyReader)
	if err != nil {
		logger.Errorf("could not create HTTP request: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf(`Bearer %s`, "eyJ2ZXJzaW9uIjoiUDJQIiwiZGF0YSI6eyJwYXlpbl9tZXJjaGFudF9zaXRlX3VpZCI6IjBldDJrMy0wMCIsInVzZXJfaWQiOiI3OTE1MzQwMDE2NSIsInNlY3JldCI6Ijc0NjQ4ZDBiZDA4YzNhYWVlZTk0NzMzMmJiZjYzODM1NmYyZWM1MmMwYjMwMGIyOTU1NDVkZjgxOTZkZTUyOWMifX0="))
	req.Header.Set("accept", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		logger.Errorf("could not send HTTP request: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			logger.Error(err)
			return
		}
	}(res.Body)

	type data struct {
		Value string `json:"value"`
	}

	type result struct {
		Status data `json:"status"`
	}

	var t result

	//t.Status.Value = "PAID"

	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		logger.Errorf("could not parse JSON response: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	logger.Debug(t)

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": user,
	})
}

// GetUserByID only one user returned
func (h *Handler) GetUserByID(ctx iris.Context) {
	rawToken := ctx.GetHeader("Authorization")
	user, err := h.CheckAuth(rawToken)
	if err != nil {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
			"data": nil,
		})
		return
	}

	if user.ID == 0 {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "User not exist",
			},
			"data": nil,
		})
		return
	}

	logger.Info(user)

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": user,
	})
}
