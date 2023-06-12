package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ivahaev/go-logger"
	"net/http"
	"os"
)

const (
	directoryPath = `./uploads/tg_chats_photos/%s`
)

func (h *Handler) CallBackVK(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	user, err := h.CheckAuth(state)
	if err != nil {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	if user == nil {
		ctx.JSON(405,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	err = h.Service.Storage.CallBackVK(code, user.Token)
	if err != nil {
		logger.Error(err)
		ctx.Redirect(301, "https://targetboost.ru/error_auth_vk")
		return
	}
	ctx.Redirect(301, "https://targetboost.ru/tasks")
}

func (h *Handler) GetPhotoFile(ctx *gin.Context) {
	key := ctx.Query("key")
	fileBytes, err := os.ReadFile(fmt.Sprintf(directoryPath, key))
	if err != nil {
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status": gin.H{
					"message": err.Error(),
				},
				"data": nil,
			})
		return
	}

	//ctx.Header("Content-Type", "image/jpeg")
	ctx.Data(200, "image/jpeg", fileBytes)
}

//
//func (h *Handler) GetFileByKey(ctx *gin.Context) {
//	key := ctx.Query("key")
//
//	file := h.Service.Storage.GetFileByKey(key)
//	if file == nil {
//		ctx.JSON(405,
//			gin.H{
//				"status": gin.H{
//					"message": "File is not exists",
//				},
//				"data": nil,
//			})
//		return
//	}
//
//	ctx.ContentType(fmt.Sprintf(`%s/%s`, file.Type, file.Ext))
//	err := ctx.ServeFile(fmt.Sprintf(`./%s/%s.%s`, file.Path, file.Key, file.Ext))
//	if err != nil {
//		logger.Error(err)
//		ctx.StatusCode(405)
//		_ = ctx.JSON(iris.Map{
//			"status": iris.Map{
//				"code":    500,
//				"message": "Server uploaded a file with an error",
//			},
//			"data": []iris.Map{},
//		})
//		return
//	}
//}
