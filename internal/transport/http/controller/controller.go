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
)

//type Controller struct {
//	Services *services.Services
//}

func NewController(ctx context.Context, services *services.Services, bot *bot.Bot) *http.Server {
	app := gin.New()
	//app.Use(globalMiddleware)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://targetboost.ru", "https://staging.targetboost.ru"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
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
