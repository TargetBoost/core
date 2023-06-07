package handler

import (
	"github.com/kataras/iris/v12"
)

func (h *Handler) IsAuth(ctx iris.Context) {
	rawToken := ctx.GetHeader("token")
	if len(rawToken) == 0 {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "Not token required",
			},
			"data": nil,
		})
		return
	}

	_, isAuth := h.Service.Account.IsAuth(rawToken)
	if !isAuth {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "Bad token required",
			},
			"data": nil,
		})
		return
	}

	ctx.Next()
}
