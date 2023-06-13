package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ivahaev/go-logger"
)

const errorAuth = "Ошибка авторизации"

func (h *Handler) IsAuth(ctx *gin.Context) {
	rawToken := ctx.GetHeader("Authorization")
	if len(rawToken) == 0 {
		logger.Debug("rawToken == 0")
		ctx.JSON(401,
			gin.H{
				"status": gin.H{
					"message": errorAuth,
				},
				"data": nil,
			})
		return
	}

	_, isAuth := h.Service.Account.IsAuth(rawToken)
	if !isAuth {
		ctx.JSON(401,
			gin.H{
				"status": gin.H{
					"message": errorAuth,
				},
				"data": nil,
			})
		return
	}

	ctx.Next()
}
