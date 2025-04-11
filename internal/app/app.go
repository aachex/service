package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/aachex/service/internal/controller"
	"github.com/aachex/service/internal/repository/postgres"
)

type App struct {
	srv *http.Server
}

func New() *App {
	return &App{}
}

func (a *App) Start(ctx context.Context) {
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

	// Старт сервера
	a.srv = &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: mux,
	}

	a.srv.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}
