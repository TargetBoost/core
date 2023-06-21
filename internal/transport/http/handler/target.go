package handler

import (
	"core/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/ivahaev/go-logger"
	"net/http"
)

func (h *Handler) GetTargets(ctx *gin.Context) {
	rawToken := ctx.GetHeader("Authorization")
	user := h.Service.Account.GetUserByToken(rawToken)
	targets := h.Service.Target.GetTargets(user.ID)

	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": targets,
		})
}

func (h *Handler) GetTargetsToAdmin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": h.Service.Target.GetTargetsToAdmin(),
		})
}

func (h *Handler) GetTargetsToExecutors(ctx *gin.Context) {
	rawToken := ctx.GetHeader("Authorization")
	user := h.Service.Account.GetUserByToken(rawToken)
	targets := h.Service.Queue.GetTargetsToExecutor(int64(user.ID))
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": targets,
		})
	return
}

func (h *Handler) CreateTarget(ctx *gin.Context) {
	var t models.TargetService
	_ = ctx.BindJSON(&t)

	rawToken := ctx.GetHeader("Authorization")
	user, err := h.CheckAuth(rawToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	err = h.Service.Target.CreateTarget(user.ID, &t)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound,
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
	return
}

func (h *Handler) UpdateTarget(ctx *gin.Context) {
	var t models.UpdateTargetService
	_ = ctx.BindJSON(&t)

	h.Service.Target.UpdateTarget(t.ID, t.Status)
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": nil,
		})
	return
}

func (h *Handler) UpdateTargetAdvertiser(ctx *gin.Context) {
	var t models.UpdateTargetService
	_ = ctx.BindJSON(&t)

	rawToken := ctx.GetHeader("Authorization")
	_, err := h.CheckAuth(rawToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	if t.Status != 3 {
		ctx.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": "У Вас нет прав перевести кампанию в такой статус",
				},
				"data": nil,
			})
		return
	}

	h.Service.Target.UpdateTarget(t.ID, t.Status)
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": nil,
		})
	return
}

func (h *Handler) CheckTarget(ctx *gin.Context) {
	var t models.UpdateTargetService
	_ = ctx.BindJSON(&t)

	rawToken := ctx.GetHeader("Authorization")
	user := h.Service.Account.GetUserByToken(rawToken)

	chatID, cost := h.Service.Queue.GetChatID(uint(t.TID))
	userChatID := h.Service.Target.GetUserID(user.ID)

	logger.Info(chatID, userChatID)

	status, err := h.Bot.CheckMembers(chatID, userChatID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	if !status {
		ctx.JSON(http.StatusOK,
			gin.H{
				"status": gin.H{
					"message": nil,
				},
				"data": nil,
			})
		return
	}

	task := h.Service.Queue.GetTaskByID(t.ID)
	if task.Status != 1 {
		ctx.AbortWithStatusJSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": "Вы уже выполнили это задание",
				},
				"data": nil,
			})
		return
	}
	h.Service.Queue.UpdateTaskStatus(t.ID)
	h.Service.Account.UpdateUserBalance(int64(user.ID), cost)
	ctx.JSON(http.StatusNotFound,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": true,
		})
	return
}
