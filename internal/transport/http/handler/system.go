package handler

import (
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
		},
	})
}
