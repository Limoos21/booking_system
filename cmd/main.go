package main

import (
	"booking_system/cmd/providers"
	"booking_system/cmd/providers/middelware"
	"booking_system/internal/app/usecase"
	"booking_system/internal/config"
	"booking_system/internal/infrastructure/adapters/controllers"
	"booking_system/internal/infrastructure/storage"
	"log/slog"
)

// @title Booking System API
// @version 1.0
// @description This is a sample server for a booking system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {
	conf := config.NewConfig()
	log, err := providers.InitLogger(conf.LogLevel)
	if err != nil {
		panic(err)
	}
	dataBase, err := providers.NewDatabase(conf.DsnDatabase)
	if err != nil {
		log.Error("Failed to connect to database", err)
		return
	}
	jwt := middelware.NewJwt(conf.TokenBot)
	st := storage.New(log, dataBase.DataBase)
	useCase := usecase.New(st, log, conf.TokenBot, jwt)
	controller := controllers.New(log, useCase, jwt)

	httpServer := providers.NewHTTPServer(conf.GetHttpPort(), conf.LogLevel, controller)

	go func(httpServer *providers.HTTPServer, logger *slog.Logger, j *middelware.Jwt) {
		httpServer.Run(logger, j)
	}(httpServer, log, jwt)
	select {}
}
