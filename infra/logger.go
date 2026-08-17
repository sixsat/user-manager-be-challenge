package infra

import (
	"log/slog"
	"os"
)

func InitLogger() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	))
}
