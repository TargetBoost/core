package controller

import (
	"context"
	"core/internal/services"
	"core/internal/tg/bot"
	"core/internal/transport/http/router"
	"github.com/iris-contrib/middleware/cors"
	"github.com/ivahaev/go-logger"
	"github.com/kataras/iris/v12"
)

//type Controller struct {
//	Services *services.Services
//}

func NewController(ctx context.Context, services *services.Services, bot *bot.Bot) {
	app := iris.New()
	app.UseGlobal(globalMiddleware)

	iris.RegisterOnInterrupt(func() {
		err := app.Shutdown(ctx)
		if err != nil {
			logger.Error(err)
		}
	})

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		Debug:            false,
	})
	app.UseRouter(crs)

	irisRouter := router.NewRouter(app, services, bot)

	err := irisRouter.Listen(":8080", iris.WithoutInterruptHandler, iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		logger.Error(err)
	}
}

func globalMiddleware(ctx iris.Context) {
	logger.Info(ctx.Request())
}
