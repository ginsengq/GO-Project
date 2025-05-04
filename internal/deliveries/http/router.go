package http

import (
	"myproject/internal/deliveries/http/handler"
	carhandler "myproject/internal/deliveries/http/handler/car"
	orderhandler "myproject/internal/deliveries/http/handler/order"
	paymenthandler "myproject/internal/deliveries/http/handler/payment"
	userhandler "myproject/internal/deliveries/http/handler/user"
	"myproject/internal/usecases/car"
	ordercase "myproject/internal/usecases/order"
	paymentcase "myproject/internal/usecases/payment"
	usercase "myproject/internal/usecases/user"
	"myproject/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	UserUC    usercase.UseCase
	CarUC     car.CarUseCase
	OrderUC   ordercase.UseCase
	PaymentUC paymentcase.PaymentUseCase
	Logger    logger.Interface
}

func NewRouter(deps RouterDependencies) *gin.Engine {
	router := gin.Default()

	commonHandler := handler.NewCommonHandler(deps.Logger)

	userHandler := userhandler.NewHandler(deps.UserUC, deps.Logger)
	carHandler := carhandler.NewHandler(deps.CarUC, deps.Logger)
	orderHandler := orderhandler.NewHandler(deps.OrderUC, deps.Logger)
	paymentHandler := paymenthandler.NewHandler(deps.PaymentUC, deps.Logger)

	router.GET("/health", commonHandler.HealthCheck)

	api := router.Group("/api")
	{
		userRoutes := api.Group("/users")
		{
			userRoutes.POST("", userHandler.CreateUser)
			userRoutes.GET("/:id", userHandler.GetUserByID)
			userRoutes.PUT("/:id", userHandler.UpdateUser)
			userRoutes.DELETE("/:id", userHandler.DeleteUser)
			userRoutes.GET("", userHandler.ListUsers)
			userRoutes.POST("/:id/password", userHandler.ChangePassword)
			userRoutes.POST("/auth", userHandler.AuthenticateUser)
		}

		carRoutes := api.Group("/cars")
		{
			carRoutes.POST("", carHandler.CreateCar)
			carRoutes.GET("/:id", carHandler.GetCar)
			carRoutes.PUT("/:id", carHandler.UpdateCar)
			carRoutes.DELETE("/:id", carHandler.DeleteCar)
			carRoutes.GET("", carHandler.ListCars)
			carRoutes.PATCH("/:id/status", carHandler.ChangeCarStatus)
		}

		orderRoutes := api.Group("/orders")
		{
			orderRoutes.POST("", orderHandler.CreateOrder)
			orderRoutes.GET("/:id", orderHandler.GetOrder)
			orderRoutes.GET("/user/:user_id", orderHandler.GetOrdersByUserID)
			orderRoutes.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
			orderRoutes.DELETE("/:id", orderHandler.CancelOrder)
			orderRoutes.GET("", orderHandler.ListAllOrders)
		}

		paymentRoutes := api.Group("/payments")
		{
			paymentRoutes.POST("/deposit", paymentHandler.Deposit)
			paymentRoutes.POST("/transactions", paymentHandler.CreateTransaction)
			paymentRoutes.GET("/user/:user_id/transactions", paymentHandler.GetTransactionsByUser)
		}
	}

	router.NoRoute(commonHandler.NotFound)
	router.NoMethod(commonHandler.MethodNotAllowed)

	return router
}
