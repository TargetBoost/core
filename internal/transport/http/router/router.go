package router

import (
	"core/internal/services"
	"core/internal/transport/http/handler"
	"github.com/kataras/iris/v12"
)

type Router struct {
	iris *iris.Application
}

func NewRouter(iris *iris.Application, services *services.Services) *iris.Application {
	serv := handler.Handler{
		Service: services,
	}

	v1 := iris.Party("/v1")
	system := v1.Party("/system")
	service := v1.Party("/service")

	// System
	system.Handle("GET", "/health_check", serv.HealthCheck)
	system.Handle("GET", "/settings", serv.Settings)
	system.Handle("POST", "/settings", serv.SetSettings)
	system.Handle("POST", "/registration", serv.CreateUser)
	system.Handle("POST", "/auth", serv.Authorization)
	system.Handle("GET", "/is_auth", serv.IsAuth)

	// User
	service.Handle("GET", "/users", serv.GetAllUsers)
	service.Handle("GET", "/user/{id:int64}", serv.GetUserByID)

	// Target
	service.Handle("GET", "/target", serv.GetTargets)
	service.Handle("PUT", "/target", serv.UpdateTarget)
	service.Handle("GET", "/admin/target", serv.GetTargetsToAdmin)
	service.Handle("GET", "/executor/target", serv.GetTargetsToExecutors)
	service.Handle("POST", "/target", serv.CreateTarget)

	// storage
	v1.Handle("GET", "/file/{key:string}", serv.GetFileByKey)

	return iris
}
