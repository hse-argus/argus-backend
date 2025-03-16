package logger

import (
	"log/slog"
	"os"
)

var (
	logger *slog.Logger
)

func InitLogger() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func Info(msg string) {
	logger.Info(msg)
}

func Error(msg string) {
	logger.Error(msg)
}
