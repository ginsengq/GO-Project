package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	configs "myproject/internal/app/config"
	. "myproject/internal/deliveries/http"
	carrepo "myproject/internal/repositories/car"
	orderrepo "myproject/internal/repositories/order"
	paymentrepo "myproject/internal/repositories/payment"
	userrepo "myproject/internal/repositories/user"
	carservice "myproject/internal/services/car"
	orderservice "myproject/internal/services/order"
	paymentservice "myproject/internal/services/payment"
	userservice "myproject/internal/services/user"
	"myproject/pkg/logger"
)

func main() {
	cfg := configs.LoadConfig()
	if cfg == nil {
		log.Fatalf("failed to load config")
		return
	}

	appLogger := logger.New(cfg.App.LogLevel)
	appLogger.Info("app started", "port", cfg.Server.Port, "environment", cfg.App.Environment)

	dbPool, err := pgxpool.Connect(context.Background(), fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode))
	if err != nil {
		appLogger.Fatal("failed to connect to database", "error", err)
	}
	defer dbPool.Close()
	appLogger.Info("database connected")

	userRepository := userrepo.NewPostgresRepo(dbPool)
	carRepository := carrepo.NewPostgresRepo(dbPool)
	orderRepository := orderrepo.NewPostgresRepository(dbPool)
	paymentRepository := paymentrepo.NewPaymentRepository(dbPool)

	userUseCase := userservice.NewUserService(userRepository)
	carUseCase := carservice.NewService(carRepository)                                                                              // Используем сервис car
	orderUseCase := orderservice.NewService(orderRepository, carUseCase, userUseCase, paymentservice.NewService(paymentRepository)) // Добавляем зависимость от CarService
	paymentUseCase := paymentservice.NewService(paymentRepository)

	routerDeps := RouterDependencies{
		UserUC:    userUseCase,
		CarUC:     carUseCase,
		OrderUC:   orderUseCase,
		PaymentUC: paymentUseCase,
		Logger:    appLogger,
	}
	router := NewRouter(routerDeps)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("failed to start server", "error", err)
		}
	}()
	appLogger.Info("http server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("server shutdown failed", "error", err)
	}
	appLogger.Info("server stopped")
}
