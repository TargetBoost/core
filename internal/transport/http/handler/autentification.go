package handler

import (
	"github.com/kataras/iris/v12"
)

const errorAuth = "Ошибка авторизации"

func (h *Handler) IsAuth(ctx iris.Context) {
	rawToken := ctx.GetHeader("Authorization")
	if len(rawToken) == 0 {
		ctx.StatusCode(401)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": errorAuth,
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
				"message": errorAuth,
			},
			"data": nil,
		})
		return
	}

	ctx.Next()
}
