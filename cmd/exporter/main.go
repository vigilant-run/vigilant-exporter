package main

import (
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
)

type App struct {
	logger *zap.Logger
}

func NewApp(logger *zap.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (a *App) Run() int {
	a.logger.Info("Starting Vigilant Exporter")

	if len(os.Args) > 1 {
		a.logger.Error("This tool takes no arguments", zap.String("args", strings.Join(os.Args, " ")))
		return 1
	}

	a.logger.Info("Running exporter...")
	a.logger.Info("Exporter completed successfully")

	return 0
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v\n", err)
	}
	defer logger.Sync()

	app := NewApp(logger)
	os.Exit(app.Run())
}
