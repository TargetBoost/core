package handler

import (
	"core/internal/models"
	"github.com/kataras/iris/v12"
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

	_, isAuth := h.Service.Auth.IsAuth(rawToken)
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
