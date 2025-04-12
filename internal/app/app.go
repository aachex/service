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
	db  *sql.DB
}

func New() *App {
	return &App{}
}

func (app *App) Start(ctx context.Context) {
	var err error

	// Подключение к бд
	app.db, err = sql.Open("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		log.Fatal(err)
	}

	err = app.db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Repositories
	users := postgres.NewUsersRepository(app.db)

	// Handlers
	mux := http.NewServeMux()

	usersController := controller.NewUsersController(users)
	usersController.RegisterHandlers(mux)

	// Старт сервера
	app.srv = &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: mux,
	}

	app.srv.ListenAndServe()
}

func (app *App) Shutdown(ctx context.Context) error {
	app.db.Close()
	return app.srv.Shutdown(ctx)
}
