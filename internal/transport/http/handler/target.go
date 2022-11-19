package handler

import (
	"core/internal/models"
	"github.com/kataras/iris/v12"
)

func (h *Handler) GetTargets(ctx iris.Context) {
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

	targets := h.Service.Target.GetTargets(user.ID)
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": targets,
	})
	return
}

func (h *Handler) CreateTarget(ctx iris.Context) {
	var t models.TargetService
	_ = ctx.ReadJSON(&t)

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

	h.Service.Target.CreateTarget(user.ID, &t)
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": nil,
	})
	return
}
