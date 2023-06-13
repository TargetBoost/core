package router

import (
	"core/internal/services"
	"core/internal/transport/http/handler"
	"core/internal/transport/tg/bot"
	"github.com/gin-gonic/gin"
	"github.com/kataras/iris/v12"
)

type Router struct {
	iris *iris.Application
}

func NewRouter(app *gin.Engine, services *services.Services, bot *bot.Bot) *gin.Engine {
	serv := handler.Handler{
		Service: services,
		Bot:     bot,
	}

	v1 := app.Group("/v1")
	admin := v1.Group("/admin", serv.IsAdmin)
	service := v1.Group("/service", serv.IsAuth)

	// System
	admin.Handle("POST", "/settings", serv.SetSettings)
	admin.Handle("GET", "/profit", serv.GetProfit)

	// Login, Registration, GetSettings (all users permission)
	v1.Handle("POST", "/registration", serv.Registration)
	v1.Handle("POST", "/login", serv.Login)
	v1.Handle("GET", "/settings", serv.GetSettings)

	// Account
	admin.Handle("GET", "/users", serv.GetAllUsers)
	service.Handle("GET", "/user/:token", serv.GetUserByToken)

	// Pay
	service.Handle("POST", "/pay", serv.Pay)
	v1.Handle("GET", "/s/pay/:id", serv.ConfirmPay)

	admin.Handle("GET", "/task_cashes", serv.GetTaskCashesAdmin)
	admin.Handle("PUT", "/task_cashes", serv.UpdateTaskCashes)

	service.Handle("GET", "/task_cashes", serv.GetTaskCashes)
	service.Handle("POST", "/task_cashes", serv.CreateTaskCashes)

	// Target
	service.Handle("GET", "/target", serv.GetTargets)
	admin.Handle("PUT", "/target", serv.UpdateTarget)
	service.Handle("PUT", "/target", serv.UpdateTargetAdvertiser)

	admin.Handle("GET", "/target", serv.GetTargetsToAdmin)
	service.Handle("GET", "/executor/target", serv.GetTargetsToExecutors)
	service.Handle("POST", "/target", serv.CreateTarget)
	service.Handle("POST", "/check_target", serv.CheckTarget)

	// storage
	//v1.Handle("GET", "/file/{key:string}", serv.GetFileByKey)
	v1.Handle("GET", "/file_ch/:key", serv.GetPhotoFile)
	v1.Handle("GET", "/callback_vk", serv.CallBackVK)

	return app
}
