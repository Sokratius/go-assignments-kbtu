package v1

import (
	"net/http"
	"practice-7/config"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/pkg/logger"
	"practice-7/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userRoutes struct {
	t   usecase.UserInterface
	l   logger.Interface
	cfg *config.Config
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface, l logger.Interface, cfg *config.Config) {
	r := &userRoutes{t: t, l: l, cfg: cfg}

	h := handler.Group("/users")
	{
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)

		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/me", r.GetMe)
			protected.GET("/protected/hello", r.ProtectedFunc)

			admin := protected.Group("/")
			admin.Use(utils.RoleMiddleware("admin"))
			admin.PATCH("/promote/:id", r.PromoteUser)
		}
	}
}

func (r *userRoutes) RegisterUser(c *gin.Context) {
	var createUserDTO entity.CreateUserDTO

	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		badRequest(c, err)
		return
	}

	hashedPassword, err := utils.HashPassword(createUserDTO.Password)
	if err != nil {
		serverError(c, err)
		return
	}

	role := "user"
	if createUserDTO.Role != "" {
		role = createUserDTO.Role
	}

	user := entity.User{
		Username: createUserDTO.Username,
		Email:    createUserDTO.Email,
		Password: hashedPassword,
		Role:     role,
	}

	createdUser, sessionID, err := r.t.RegisterUser(&user)
	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully",
		"session_id": sessionID,
		"user":       createdUser,
	})
}

func (r *userRoutes) LoginUser(c *gin.Context) {
	var input entity.LoginUserDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		badRequest(c, err)
		return
	}

	token, err := r.t.LoginUser(&input)
	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) GetMe(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is missing in token"})
		return
	}

	userIDString, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id format in token"})
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	user, err := r.t.GetMe(userID)
	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (r *userRoutes) PromoteUser(c *gin.Context) {
	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, err)
		return
	}

	var dto entity.PromoteUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		badRequest(c, err)
		return
	}

	updatedUser, err := r.t.PromoteUser(targetID, dto.Role)
	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user role updated",
		"user":    updatedUser,
	})
}

func (r *userRoutes) ProtectedFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
