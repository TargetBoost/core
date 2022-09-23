package handler

import (
	"core/internal/models"
	"github.com/kataras/iris/v12"
	"time"
)

func (h *Handler) HealthCheck(ctx iris.Context) {
	ctx.StatusCode(200)

	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": iris.Map{
			"time": time.Now().Unix(),
		},
	})
}

func (h *Handler) Settings(ctx iris.Context) {
	ctx.StatusCode(200)

	settings := h.Service.Settings.GetSettings()

	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": iris.Map{
			"snow": settings.Snow,
			"rain": settings.Rain,
		},
	})
}

func (h *Handler) SetSettings(ctx iris.Context) {
	var s models.Settings
	err := ctx.ReadJSON(&s)
	if err != nil {
		ctx.StatusCode(400)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": "bad data insertion",
			},
			"data": nil,
		})
		return
	}

	rawToken := ctx.GetHeader("Authorization")
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

	uid, isAuth := h.Service.Auth.IsAuth(rawToken)
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

	user := h.Service.User.GetUserByID(int64(uid))
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

	h.Service.Settings.SetSettings(&s)
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": nil,
	})
}
