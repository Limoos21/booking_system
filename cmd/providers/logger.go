package providers

import (
	"errors"
	"log/slog"
	"os"
)

const (
	LevelDebug   = "debug"
	LevelRelease = "release"
	LevelTest    = "test"
)

func InitLogger(logLevel string) (*slog.Logger, error) {
	var level slog.Level

	switch logLevel {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelRelease:
		level = slog.LevelInfo
	case LevelTest:
		level = slog.LevelWarn
	default:
		return nil, errors.New("invalid log level. Available levels: debug, release, test")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	return logger, nil
}
