package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	app := app.New()
	go app.Start()

	// shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interrupt
	fmt.Println("shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = app.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatal(err)
	}
}
