package handler

import (
	"bytes"
	"core/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ivahaev/go-logger"
	"github.com/kataras/iris/v12"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Registration - registration user
func (h *Handler) Registration(ctx *gin.Context) {
	// userData - data new user
	var userData models.CreateUser
	err := ctx.BindJSON(&userData)
	if err != nil {
		ctx.JSON(400,
			gin.H{
				"status": gin.H{
					"message": errorDataInsertion,
				},
				"data": nil,
			})
		return
	}

	// Registration - service returned data for new user
	user, err := h.Service.Account.CreateUser(userData)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": fmt.Sprintf(errorService, "регистрации"),
				},
				"data": nil,
			})
		return
	}
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": gin.H{
				"token": user.Token,
				"id":    user.ID,
			},
		})
}

// GetAllUsers all users returned for admin
func (h *Handler) GetAllUsers(ctx *gin.Context) {
	users := h.Service.Account.GetAllUsers()
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": users,
		})
}

// CreateTaskCashes create
func (h *Handler) CreateTaskCashes(ctx *gin.Context) {
	rawToken := ctx.GetHeader("Authorization")
	var task models.TaskCashToUser
	err := ctx.BindJSON(&task)
	if err != nil {
		ctx.JSON(400,
			gin.H{
				"status": gin.H{
					"message": "bad data insertion",
				},
				"data": nil,
			})
		return
	}

	u := h.Service.Account.GetUserByToken(rawToken)
	err = h.Service.Account.CreateTaskCashes(int64(u.ID), task)
	if err != nil {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": nil,
		})
}

// UpdateTaskCashes create
func (h *Handler) UpdateTaskCashes(ctx *gin.Context) {
	var task models.TaskCashToService
	err := ctx.BindJSON(&task)
	if err != nil {
		ctx.JSON(400,
			gin.H{
				"status": gin.H{
					"message": "bad data insertion",
				},
				"data": nil,
			})
		return
	}

	h.Service.Account.UpdateTaskCashes(task)
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": nil,
		})
}

func (h *Handler) GetTaskCashes(ctx *gin.Context) {
	rawToken := ctx.GetHeader("Authorization")
	u := h.Service.Account.GetUserByToken(rawToken)
	tasks := h.Service.Account.GetTasksCashesUser(u.ID)
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": tasks,
		})
}

func (h *Handler) GetTaskCashesAdmin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": h.Service.Account.GetTasksCashesAdmin(),
		})
}

// Pay -
func (h *Handler) Pay(ctx *gin.Context) {
	rawToken := ctx.GetHeader("Authorization")
	var pay models.Pay
	_ = ctx.BindJSON(&pay)

	user := h.Service.Account.GetUserByToken(rawToken)

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

	jsonBody, err := json.Marshal(body)
	if err != nil {
		logger.Error(err)
	}

	httpClient := http.Client{}

	reqURL := fmt.Sprintf("https://api.qiwi.com/partner/bill/v1/bills/%s", id.String())
	req, err := http.NewRequest(http.MethodPut, reqURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Errorf("could not create HTTP request: %v", err)
		ctx.Status(iris.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf(`Bearer %s`, "eyJ2ZXJzaW9uIjoiUDJQIiwiZGF0YSI6eyJwYXlpbl9tZXJjaGFudF9zaXRlX3VpZCI6IjBldDJrMy0wMCIsInVzZXJfaWQiOiI3OTE1MzQwMDE2NSIsInNlY3JldCI6Ijc0NjQ4ZDBiZDA4YzNhYWVlZTk0NzMzMmJiZjYzODM1NmYyZWM1MmMwYjMwMGIyOTU1NDVkZjgxOTZkZTUyOWMifX0="))
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	res, err := httpClient.Do(req)
	if err != nil {
		logger.Errorf("could not send HTTP request: %v", err)
		ctx.Status(iris.StatusInternalServerError)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ctx.Status(iris.StatusInternalServerError)
			logger.Error(err)
			return
		}
	}(res.Body)

	type Result struct {
		SiteID string `json:"siteId"`
		BillID string `json:"billId"`
		Amount struct {
			Currency string `json:"currency"`
			Value    string `json:"value"`
		} `json:"amount"`
		Status struct {
			Value           string    `json:"value"`
			ChangedDateTime time.Time `json:"changedDateTime"`
		} `json:"status"`
		CreationDateTime     time.Time `json:"creationDateTime"`
		ExpirationDateTime   time.Time `json:"expirationDateTime"`
		PayURL               string    `json:"payUrl"`
		RecipientPhoneNumber string    `json:"recipientPhoneNumber"`
	}

	var t Result

	//t.Status.Value = "PAID"

	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		logger.Errorf("could not parse JSON response: %v", err)
		ctx.Status(iris.StatusInternalServerError)
		return
	}

	if t.Status.Value != "WAITING" {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": "Ошибка шлюза",
				},
				"data": nil,
			})
		return
	}

	var trans models.TransactionToService
	trans.BuildID = id.String()
	trans.UID = user.ID
	trans.Status = t.Status.Value
	trans.Amount = pay.Value

	h.Service.Account.CreateTransaction(&trans)
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": gin.H{
				"url": fmt.Sprintf("https://oplata.qiwi.com/create?publicKey=48e7qUxn9T7RyYE1MVZswX1FRSbE6iyCj2gCRwwF3Dnh5XrasNTx3BGPiMsyXQFNKQhvukniQG8RTVhYm3iP3f4HArt65TUfZCPMYWpVH2XN4KRVBdZrHB6RTHkUcsdeHGekuM4JXb4Cd5JvDucawYX8bSof9fjuacyrjAfPGRNegJXbgdK19u2QSSwVk&billId=%s&amount=%s&account=5&customFields[themeCode]=Andrei-ShQU6cQ2pop&successUrl=https://targetboost.ru/core/v1/s/pay/%s", id.String(), pay.Value, id.String()),
			},
		})
	return
}

func (h *Handler) ConfirmPay(ctx *gin.Context) {
	key := ctx.Query("id")
	trans := h.Service.Account.GetTransaction(key)
	if trans.Status == "" {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": "Transaction is not validate",
				},
				"data": nil,
			})
		return
	}

	if trans.Status != "WAITING" {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": "Transaction is PAID before",
				},
				"data": nil,
			})
		return
	}

	httpClient := http.Client{}

	reqURL := fmt.Sprintf("https://api.qiwi.com/partner/bill/v1/bills/%s", trans.BuildID)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		logger.Errorf("could not create HTTP request: %v", err)
		ctx.Status(iris.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf(`Bearer %s`, "eyJ2ZXJzaW9uIjoiUDJQIiwiZGF0YSI6eyJwYXlpbl9tZXJjaGFudF9zaXRlX3VpZCI6IjBldDJrMy0wMCIsInVzZXJfaWQiOiI3OTE1MzQwMDE2NSIsInNlY3JldCI6Ijc0NjQ4ZDBiZDA4YzNhYWVlZTk0NzMzMmJiZjYzODM1NmYyZWM1MmMwYjMwMGIyOTU1NDVkZjgxOTZkZTUyOWMifX0="))
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	res, err := httpClient.Do(req)
	if err != nil {
		logger.Errorf("could not send HTTP request: %v", err)
		ctx.Status(iris.StatusInternalServerError)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ctx.Status(iris.StatusInternalServerError)
			logger.Error(err)
			return
		}
	}(res.Body)

	type Result struct {
		SiteID string `json:"siteId"`
		BillID string `json:"billId"`
		Amount struct {
			Currency string `json:"currency"`
			Value    string `json:"value"`
		} `json:"amount"`
		Status struct {
			Value           string    `json:"value"`
			ChangedDateTime time.Time `json:"changedDateTime"`
		} `json:"status"`
		CreationDateTime     time.Time `json:"creationDateTime"`
		ExpirationDateTime   time.Time `json:"expirationDateTime"`
		PayURL               string    `json:"payUrl"`
		RecipientPhoneNumber string    `json:"recipientPhoneNumber"`
	}

	var t Result

	//t.Status.Value = "PAID"

	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		logger.Errorf("could not parse JSON response: %v", err)
		ctx.Status(iris.StatusInternalServerError)
		return
	}

	if t.Status.Value != "PAID" {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": "Transaction is WAIT status",
				},
				"data": nil,
			})
		return
	}

	f, err := strconv.ParseFloat(trans.Amount, 64)
	if err != nil {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	u := h.Service.Account.GetUserByID(int64(trans.UID))

	fu, err := strconv.ParseFloat(u.Balance, 64)
	if err != nil {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	h.Service.Account.UpdateUser(u.ID, fu+f)
	trans.Status = "PAID"
	h.Service.Account.UpdateTransaction(trans)
	ctx.Redirect(301, "https://targetboost.ru/s/pay")
}

// GetUserByToken only one user returned
func (h *Handler) GetUserByToken(ctx *gin.Context) {
	rawToken := ctx.GetHeader("Authorization")
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": h.Service.Account.GetUserByToken(rawToken),
		})
}

// isAdmin check if admin middleware
func (h *Handler) IsAdmin(ctx *gin.Context) {
	rawToken := ctx.GetHeader("Authorization")
	user, err := h.CheckAuth(rawToken)
	if err != nil {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	if user.Admin != true {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": "У Вас нет прав доступа",
				},
				"data": nil,
			})
		return
	}

	ctx.Next()
}
