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
				"message": nil,
			},
			"data": nil,
		})
	}

	h.Service.User.CreateUser(u)

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
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
