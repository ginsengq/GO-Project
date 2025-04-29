package auth

import (
	service "myproject/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewAuthHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) signUp(c *gin.Context) {
	var userData struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.services.AuthService.Register(userData.Name, userData.Email, userData.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (h *Handler) signIn(c *gin.Context) {
	var userData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Проверка данных и создание токена
	token, err := h.services.AuthService.Login(userData.Email, userData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func RegisterRoutes(router *gin.Engine, services *service.Service) {
	h := NewAuthHandler(services)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
}
