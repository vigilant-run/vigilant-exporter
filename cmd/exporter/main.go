package main

import (
	"log/slog"
	"os"
	"strings"
)

type App struct {
	logger *slog.Logger
}

func NewApp(
	logger *slog.Logger,
) *App {
	return &App{
		logger: logger,
	}
}

func (a *App) Run() int {
	if len(os.Args) > 1 {
		a.logger.Error("This tool takes no arguments", slog.String("args", strings.Join(os.Args, " ")))
		return 1
	}

	return 0
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := NewApp(logger)
	os.Exit(app.Run())
}
