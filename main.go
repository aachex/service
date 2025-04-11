package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/aachex/service/internal/app"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	bgCtx := context.Background()

	ctx, stop := signal.NotifyContext(bgCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := app.New()
	go app.Start(ctx)

	<-ctx.Done()

	fmt.Println("shutdown")

	shutdownCtx, cancel := context.WithTimeout(bgCtx, 5*time.Second)
	defer cancel()

	err = app.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatal(err)
	}
}
