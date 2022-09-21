package handler

import (
	"core/internal/models"
	"github.com/kataras/iris/v12"
)

func (h *Handler) Authorization(ctx iris.Context) {
	var a models.AuthUser
	err := ctx.ReadJSON(&a)
	if err != nil {
		ctx.StatusCode(400)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": nil,
			},
			"data": nil,
		})
		return
	}

	token, err := h.Service.User.AuthUser(a)
	if err != nil {
		ctx.StatusCode(400)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err,
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
		"data": iris.Map{
			"token": token,
		},
	})
	return
}
