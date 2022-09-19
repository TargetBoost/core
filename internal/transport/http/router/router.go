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
	service := v1.Party("/service", serv.AuthMiddleware)

	// System
	system.Handle("GET", "/health_check", serv.HealthCheck)
	system.Handle("GET", "registration", serv.)

	// User
	service.Handle("GET", "/users", serv.GetAllUsers)
	service.Handle("GET", "/user/{id:int64}", serv.GetUserByID)

	// feed
	v1.Handle("GET", "/feeds", serv.GetAllFeeds)
	v1.Handle("GET", "/feed/{id:int64}", serv.GetFeedByID)

	// storage
	v1.Handle("GET", "/file/{key:string}", serv.GetFileByKey)

	return iris
}
