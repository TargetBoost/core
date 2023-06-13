package handler

import (
	"core/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) Login(ctx *gin.Context) {
	var a models.AuthUser
	err := ctx.BindJSON(&a)
	if err != nil {
		ctx.AbortWithStatusJSON(400,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	user, err := h.Service.Account.AuthUser(a)
	if err != nil {
		ctx.AbortWithStatusJSON(400,
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
			"data": gin.H{
				"token":   user.Token,
				"id":      user.ID,
				"execute": user.Execute,
			},
		})
}
