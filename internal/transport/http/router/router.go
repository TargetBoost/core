package router

import (
	"core/internal/services"
	"core/internal/transport/http/handler"
	"core/internal/transport/tg/bot"
	"github.com/kataras/iris/v12"
)

type Router struct {
	iris *iris.Application
}

func NewRouter(iris *iris.Application, services *services.Services, bot *bot.Bot) *iris.Application {
	serv := handler.Handler{
		Service: services,
		Bot:     bot,
	}

	v1 := iris.Party("/v1")
	admin := v1.Party("/admin", serv.IsAdmin)
	service := v1.Party("/service", serv.IsAuth)

	// System
	admin.Handle("GET", "/settings", serv.Settings)
	admin.Handle("POST", "/settings", serv.SetSettings)

	// Login and Registration
	v1.Handle("POST", "/registration", serv.Registration)
	v1.Handle("POST", "/login", serv.Login)

	//service.Handle("GET", "/is_auth", serv.IsAuth)

	// Account
	service.Handle("GET", "/users", serv.GetAllUsers)
	service.Handle("GET", "/user/{id:int64}", serv.GetUserByID)
	service.Handle("POST", "/pay", serv.Pay)
	service.Handle("GET", "/s/pay/{id:string}", serv.ConfirmPay)

	service.Handle("GET", "/admin/task_cashes", serv.GetTaskCashesAdmin)
	service.Handle("PUT", "/admin/task_cashes", serv.UpdateTaskCashes)

	service.Handle("GET", "/task_cashes", serv.GetTaskCashes)
	service.Handle("POST", "/task_cashes", serv.CreateTaskCashes)

	// Target
	service.Handle("GET", "/target", serv.GetTargets)
	service.Handle("PUT", "/admin/target", serv.UpdateTarget)
	service.Handle("PUT", "/target", serv.UpdateTargetAdvertiser)

	service.Handle("GET", "/admin/target", serv.GetTargetsToAdmin)
	service.Handle("GET", "/executor/target", serv.GetTargetsToExecutors)
	service.Handle("POST", "/target", serv.CreateTarget)
	service.Handle("POST", "/check_target", serv.CheckTarget)

	// storage
	v1.Handle("GET", "/file/{key:string}", serv.GetFileByKey)
	v1.Handle("GET", "/file_ch/{key:string}", serv.GetPhotoFile)
	v1.Handle("GET", "/callback_vk", serv.CallBackVK)

	return iris
}
