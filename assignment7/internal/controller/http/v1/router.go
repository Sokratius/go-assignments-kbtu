package v1

import (
	"practice-7/config"
	"practice-7/internal/usecase"
	"practice-7/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.RouterGroup, t usecase.UserInterface, l logger.Interface, cfg *config.Config) {
	newUserRoutes(handler, t, l, cfg)
}
