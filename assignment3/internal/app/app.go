package app

import (
	"golang/internal/delivery/http"
	"golang/internal/repository/_postgres"
	"golang/internal/repository/_postgres/users"
	"golang/internal/usecase"
)

type App struct {
	Handler *http.Handler
}

func NewApp(db *_postgres.Dialect) *App {
	userRepo := users.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	handler := http.NewHandler(userUsecase)

	return &App{
		Handler: handler,
	}
}
