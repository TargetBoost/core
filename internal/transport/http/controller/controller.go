package controller

import (
	"context"
	"core/internal/services"
	"core/internal/transport/http/router"
	"core/internal/transport/tg/bot"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ivahaev/go-logger"
	"net/http"
	"time"
)

//type Controller struct {
//	Services *services.Services
//}

func NewController(ctx context.Context, services *services.Services, bot *bot.Bot) *http.Server {
	app := gin.Default()
	app.Use(globalMiddleware)

	//iris.RegisterOnInterrupt(func() {
	//	err := app.Shutdown(ctx)
	//	if err != nil {
	//		logger.Error(err)
	//	}
	//})

	//crs := cors.New(cors.Options{
	//	AllowedOrigins:   []string{"*"},
	//	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	//	AllowCredentials: true,
	//	Debug:            false,
	//})
	//app.Use(crs)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://targetboost.ru", "https://staging.targetboost.ru"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://targetboost.ru"
		},
		MaxAge: 12 * time.Hour,
	}))

	irisRouter := router.NewRouter(app, services, bot)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: irisRouter,
	}

	return srv
}

func globalMiddleware(ctx *gin.Context) {
	logger.Info(
		fmt.Sprintf("URL: %s", ctx.Request.URL),
		fmt.Sprintf("Header: %s", ctx.Request.Header),
		fmt.Sprintf("Method: %s", ctx.Request.Method),
		fmt.Sprintf("Host: %s", ctx.Request.Host),
	)
	ctx.Next()
}
