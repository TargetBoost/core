package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetBlog(ctx *gin.Context) {
	blog := h.Service.Blog.GetBlog()

	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": blog,
		})
	return
}
