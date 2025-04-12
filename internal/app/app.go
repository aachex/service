package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/aachex/service/internal/controller"
	"github.com/aachex/service/internal/repository/postgres"
)

type App struct {
	srv    *http.Server
	db     *sql.DB
	logger *slog.Logger
}

func New(l *slog.Logger) *App {
	return &App{logger: l}
}

func (app *App) Start() {
	var err error

	// Подключение к бд
	app.db, err = sql.Open("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
	err = app.db.Ping()
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
	app.logger.Info("connected to db")

	// Репозитории
	users := postgres.NewUsersRepository(app.db)

	// Обаботчики
	mux := http.NewServeMux()

	usersController := controller.NewUsersController(users, app.logger)
	usersController.RegisterHandlers(mux)

	// Старт сервера
	app.srv = &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: mux,
	}

	app.srv.ListenAndServe()
}

func (app *App) Shutdown(ctx context.Context) error {
	err := app.srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = app.db.Close()
	if err != nil {
		return err
	}

	app.logger.Info("shutdown")

	return nil
}
