package app

import (
	"fmt"
	"practice-7/config"
	v1 "practice-7/internal/controller/http/v1"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/logger"
	"practice-7/pkg/postgres"
	"practice-7/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	cfg    *config.Config
	router *gin.Engine
}

func New(cfg *config.Config) (*App, error) {
	l := logger.New()

	pg, err := postgres.New(cfg.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	userRepo, err := repo.NewUserRepo(pg)
	if err != nil {
		return nil, fmt.Errorf("init user repo: %w", err)
	}

	userUseCase := usecase.NewUserUseCase(userRepo, cfg.JWTSecret)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(utils.RateLimiterMiddleware(
		cfg.JWTSecret,
		cfg.RateLimitRequests,
		time.Duration(cfg.RateLimitWindowS)*time.Second,
	))

	apiV1 := r.Group("/api/v1")
	v1.NewRouter(apiV1, userUseCase, l, cfg)

	return &App{cfg: cfg, router: r}, nil
}

func (a *App) Run() error {
	return a.router.Run(":" + a.cfg.HTTPPort)
}
