package logger

import (
	"log/slog"
	"os"
)

var (
	log *slog.Logger
)

func InitLogger() {
	log = slog.New(slog.NewTextHandler(os.Stdin, nil))
}

func Info(info string, args ...any) {
	log.Info(info, args)
}

func Error(info string, args ...any) {
	log.Error(info, args)
}
