package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/aachex/service/internal/controller"
	"github.com/aachex/service/internal/repository/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	// Подключение к бд
	db, err := sql.Open("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	users := postgres.NewUsersRepository(db)

	// Handlers
	mux := http.NewServeMux()

	usersController := controller.NewUsersController(users)
	usersController.RegisterHandlers(mux)

	port := os.Getenv("PORT")
	addr := ":" + port
	http.ListenAndServe(addr, mux)
}
