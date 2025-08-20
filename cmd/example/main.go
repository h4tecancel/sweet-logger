package main

import (
	"errors"
	"log/slog"
	"time"

	"github.com/h4tecancel/sweet-logger"
)

func main() {
	logger := sweetlogger.New(sweetlogger.Options{
		Level:      slog.LevelDebug,
		AddSource:  true,
		TimeFormat: "15:04:05.000",
		Color:      sweetlogger.ColorAuto, // Auto/Always/Never
	})

	logger = logger.With("app", "sweet-logger", "env", "local")

	logger.Debug("debug line", "tries", 3)
	logger.Info("server started", "addr", "127.0.0.1:8080", "mode", "local")
	logger.Warn("cache miss", "key", "user:42")
	logger.Error("db error", "err", errors.New("connection refused"))

	logger.WithGroup("http").Info("request",
		"method", "GET",
		"path", "/api/v1/ping",
		"status", 200,
		"dur", 12*time.Millisecond,
	)
}
