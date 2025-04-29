package handler

import (
	service "myproject/internal/services"
	car.go "myproject/internal/deliveries/car"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api")
	{
		cars := api.Group("/cars")
		{
			cars.POST("/", h.createCar)
			cars.GET("/", h.getAllCars)
			cars.GET("/:id", h.getCarById)
			cars.PUT("/:id", h.updateCar)
			cars.DELETE("/:id", h.deleteCar)
		}

		orders := api.Group("/orders")
		{
			orders.POST("/", h.createOrder)
			orders.GET("/", h.getAllOrders)
			orders.GET("/:id", h.getOrderById)
			orders.DELETE("/:id", h.deleteOrder)
		}

		payments := api.Group(("/payments"))
		{
			payments.POST("/deposit", h.deposit)
			payments.GET("/transactions", h.getTransactions)
		}
	}

	return router
}
