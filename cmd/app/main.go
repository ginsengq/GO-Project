package main

import (
	"fmt"
	"log"
	handler "myproject/internal/delivery/http"
	"myproject/internal/repository"
	"myproject/internal/server"
	"myproject/internal/service"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("Hello!")

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("server.port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while runnng http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
