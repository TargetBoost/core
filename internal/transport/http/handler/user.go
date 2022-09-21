package handler

import (
	"core/internal/models"
	"github.com/kataras/iris/v12"
	"strings"
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

// GetAllUsers all users returned
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

	token, err := h.Service.User.CreateUser(u)
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
			"token": token,
		},
	})
}

// GetUserByID only one user returned
func (h *Handler) GetUserByID(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("id")
	if err != nil {
		ctx.StatusCode(400)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "Insert id is not int",
			},
			"data": nil,
		})
		return
	}

	rawToken := ctx.GetHeader("Authorization")
	if len(rawToken) == 0 {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "Not token required",
			},
			"data": nil,
		})
		return
	}

	contain := strings.Contains(rawToken, "Bearer")
	if !contain {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "Insert token is not validate",
			},
			"data": nil,
		})
		return
	}

	sliceToken := strings.Split(rawToken, " ")
	if len(sliceToken) == 1 || len(sliceToken) > 2 {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "Insert token is not validate",
			},
			"data": nil,
		})
		return
	}

	isAuth := h.Service.Auth.IsAuth(sliceToken[1])
	if !isAuth {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "Bad token required",
			},
			"data": nil,
		})
		return
	}

	user := h.Service.User.GetUserByID(id)
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

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": user,
	})
}
