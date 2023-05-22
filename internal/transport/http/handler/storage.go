package handler

import (
	"fmt"
	"github.com/ivahaev/go-logger"
	"github.com/kataras/iris/v12"
)

func (h *Handler) GetFileByKey(ctx iris.Context) {
	key := ctx.Params().GetString("key")

	file := h.Service.Storage.GetFileByKey(key)
	if file == nil {
		ctx.StatusCode(405)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"code":    405,
				"message": "File is not exists",
			},
			"data": []iris.Map{},
		})
		return
	}

	ctx.ContentType(fmt.Sprintf(`%s/%s`, file.Type, file.Ext))
	err := ctx.ServeFile(fmt.Sprintf(`./%s/%s.%s`, file.Path, file.Key, file.Ext))
	if err != nil {
		logger.Error(err)
		ctx.StatusCode(405)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"code":    500,
				"message": "Server uploaded a file with an error",
			},
			"data": []iris.Map{},
		})
		return
	}
}

func (h *Handler) TestVast(ctx iris.Context) {
	ctx.StatusCode(200)
	_ = ctx.JSON(iris.Map{
		"status": iris.Map{
			"code":    500,
			"message": "Server uploaded a file with an error",
		},
		"data": ``,
	})
	return
}
