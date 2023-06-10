package handler

import (
	"core/internal/models"
	"github.com/kataras/iris/v12"
)

func (h *Handler) GetSettings(ctx iris.Context) {
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

	h.Service.Settings.SetSettings(&s)
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"message": nil,
		},
		"data": nil,
	})
}
