package handler

import (
	"core/internal/models"
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

func (h *Handler) CreateComment(ctx *gin.Context) {
	rawToken := ctx.GetHeader("Authorization")
	var comment models.Comment

	err := ctx.BindJSON(&comment)
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

	h.Service.Blog.AddComment(comment, rawToken)
	ctx.JSON(http.StatusOK,
		gin.H{
			"status": gin.H{
				"message": nil,
			},
			"data": nil,
		})
	return
}
