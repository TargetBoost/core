package handler

import (
	"core/internal/models"
	"github.com/ivahaev/go-logger"
	"github.com/kataras/iris/v12"
)

func (h *Handler) GetTargets(ctx iris.Context) {
	rawToken := ctx.GetHeader("Authorization")
	user := h.Service.Account.GetUserByToken(rawToken)
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
	rawToken := ctx.GetHeader("Login")
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
	rawToken := ctx.GetHeader("Login")
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

	targets := h.Service.Queue.GetTargetsToExecutor(int64(user.ID))
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

	rawToken := ctx.GetHeader("Login")
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

	rawToken := ctx.GetHeader("Login")
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

	h.Service.Target.UpdateTarget(t.ID, t.Status)
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": nil,
	})
	return
}

func (h *Handler) UpdateTargetAdvertiser(ctx iris.Context) {
	var t models.UpdateTargetService
	_ = ctx.ReadJSON(&t)

	rawToken := ctx.GetHeader("Login")
	_, err := h.CheckAuth(rawToken)
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

	if t.Status != 3 {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": `У Вас нет прав перевести кампанию в такой статус`,
			},
			"data": nil,
		})
		return
	}

	h.Service.Target.UpdateTarget(t.ID, t.Status)
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": nil,
	})
	return
}

func (h *Handler) CheckTarget(ctx iris.Context) {
	var t models.UpdateTargetService
	_ = ctx.ReadJSON(&t)

	rawToken := ctx.GetHeader("Login")
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

	chatID, cost := h.Service.Queue.GetChatID(uint(t.TID))
	userChatID := h.Service.Target.GetUserID(user.ID)

	logger.Info(chatID, userChatID)

	status, err := h.Bot.CheckMembers(chatID, userChatID)
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

	if !status {
		ctx.StatusCode(200)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": nil,
			},
			"data": false,
		})
		return
	}

	task := h.Service.Queue.GetTaskByID(t.ID)
	if task.Status != 1 {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "Вы уже выполнили это задание",
			},
			"data": nil,
		})
		return
	}
	h.Service.Queue.UpdateTaskStatus(t.ID)
	h.Service.Account.UpdateUserBalance(int64(user.ID), cost)

	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": true,
	})
	return
}
