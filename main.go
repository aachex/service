package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aachex/service/internal/app"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver
)

func main() {
	// Логгер
	logFile, err := os.OpenFile("app.log", os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(logFile, logOpts))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		logger.Error(err.Error())
	}

	// Запуск приложения
	app := app.New(logger)
	go app.Start()

	// shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interrupt
	fmt.Println("shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		logger.Error(err.Error())
	}
}
