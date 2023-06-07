package main

import (
	"core/internal/repositories"
	"core/internal/services"
	"core/internal/target_broker"
	"core/internal/transport/http/controller"
	"core/internal/transport/tg/bot"

	"context"
	"fmt"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivahaev/go-logger"
	"gorm.io/driver/postgres"
)

const (
	shutDownDuration = 5 * time.Second
)

func main() {
	err := logger.SetLevel("debug")
	if err != nil {
		panic(fmt.Sprintf(`failed init logs, error: %s`, err.Error()))
	}

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow`,
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_PORT"),
	)

	logger.Debug(dsn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Notice("App Run ...")

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	repo := repositories.NewRepositories(db)

	b, err := bot.New(ctx, "5911800604:AAFN65f8vQrsgjIxR8vQgUr_SBCj8SQ1RoM", repo)
	if err != nil {
		panic(err)
	}

	q := target_broker.New(ctx, repo, b)
	go q.Broker()
	go q.AppointTask()
	go q.AntiFraud()

	go b.GetUpdates()
	go b.SenderUpdates()

	serv := services.NewServices(repo, q.Line, q.LineAppoint, b.TrackMessages)
	go controller.NewController(ctx, serv, b)

	<-GracefulShutdown()
	_, forceCancel := context.WithTimeout(ctx, shutDownDuration)

	logger.Notice("Graceful Shutdown")

	defer forceCancel()
}

func GracefulShutdown() chan os.Signal {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	return done
}
