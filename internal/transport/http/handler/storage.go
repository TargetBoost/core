package handler

import (
	"fmt"
	"github.com/ivahaev/go-logger"
	"github.com/kataras/iris/v12"
	"os"
)

const (
	directoryPath = `./uploads/tg_chats_photos/%s`
)

func (h *Handler) CallBackVK(ctx iris.Context) {
	code := ctx.Params().GetString("code")

	rawToken := ctx.GetHeader("Authorization")
	user, err := h.CheckAuth(rawToken)
	if err != nil {
		ctx.StatusCode(404)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
			"data": nil,
		})
		return
	}

	if user == nil {
		ctx.StatusCode(405)
		_ = ctx.JSON(iris.Map{
			"status": iris.Map{
				"message": err.Error(),
			},
			"data": nil,
		})
		return
	}

	err = h.Service.Storage.CallBackVK(code, user.Token)
	if err != nil {
		logger.Error(err)
		ctx.Redirect("https://targetboost.ru/error_auth_vk", 301)
		return
	}
	ctx.Redirect("https://targetboost.ru/tasks", 301)
}

func (h *Handler) GetPhotoFile(ctx iris.Context) {
	key := ctx.Params().GetString("key")
	fileBytes, err := os.ReadFile(fmt.Sprintf(directoryPath, key))
	if err != nil {
		ctx.StatusCode(404)
		return
	}

	ctx.StatusCode(200)
	ctx.ContentType("image/jpeg")
	ctx.Write(fileBytes)
}

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
