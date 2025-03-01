package main

import (
	"back_api/internal/config"
	"back_api/internal/repository"
	"back_api/internal/repository/models"
	"log/slog"
	"os"
)

func main() {
	config := config.NewConfig()
	dsn := config.GetDSN()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	newRepository := repository.NewPostgresRepository(dsn, logger)
	models.Migrate(newRepository.Database, logger)
	return

}
