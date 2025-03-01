package main

import (
	"back_api/internal/config"
	"back_api/internal/controllers/user_controllers"
	"back_api/internal/interfaces/api_server"
	"back_api/internal/interfaces/middelware"
	"back_api/internal/repository"
	"back_api/internal/service/user_service"
	"log/slog"
	"os"
)

const (
	GLOBAL_PREFIX = "api/v1"
)

func main() {
	config := config.NewConfig()
	dsn := config.GetDSN()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	newJWT := middelware.NewJwt(config.TokenBot)

	newRepository := repository.NewPostgresRepository(dsn, logger)
	newUserService := user_service.NewUserService(logger, newRepository, config.TokenBot)
	newServer := api_server.NewServer("localhost:6655")
	group := newServer.SetGlobalPrefix(GLOBAL_PREFIX)
	NewController := user_controllers.NewUserController(logger, newUserService, newJWT)
	api_server.NewApiRoutersUser(group, logger, NewController, newJWT)
	err := newServer.Run()
	if err != nil {
		panic(err)
	}
}
