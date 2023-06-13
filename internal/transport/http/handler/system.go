package handler

import (
	"core/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetSettings(ctx *gin.Context) {
	settings := h.Service.Settings.GetSettings()
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": gin.H{
				"snow": settings.Snow,
				"rain": settings.Rain,
			},
		})
}

func (h *Handler) GetProfit(ctx *gin.Context) {
	profit := h.Service.Target.GetProfit()
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": gin.H{
				"profit": profit,
			},
		})
}

func (h *Handler) SetSettings(ctx *gin.Context) {
	var s models.Settings
	err := ctx.BindJSON(&s)
	if err != nil {
		ctx.AbortWithStatusJSON(400,
			gin.H{
				"status": gin.H{
					"message": "bad data insertion",
				},
				"data": nil,
			})
		return
	}

	h.Service.Settings.SetSettings(&s)
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": nil,
		})
}
