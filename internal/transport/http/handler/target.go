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

func (h *Handler) GetTargetsToAdmin(ctx iris.Context) {
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

	if !user.Admin {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "your dont have permission",
			},
			"data": nil,
		})
		return
	}

	targets := h.Service.Target.GetTargetsToAdmin()
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": targets,
	})
	return
}

func (h *Handler) GetTargetsToExecutors(ctx iris.Context) {
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

	targets := h.Service.Target.GetTargetsToExecutor(int64(user.ID))
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

	err = h.Service.Target.CreateTarget(user.ID, &t)
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

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": nil,
	})
	return
}

func (h *Handler) UpdateTarget(ctx iris.Context) {
	var t models.UpdateTargetService
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

	if !user.Admin {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "your dont have permission",
			},
			"data": nil,
		})
		return
	}

	h.Service.Target.UpdateTarget(t.ID)
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": nil,
	})
	return
}
