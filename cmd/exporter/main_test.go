package main

import (
	"os"
	"testing"

	"log/slog"
)

func TestRunWithNoArgs(t *testing.T) {
	logger := newTestLogger(t)
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"vigilant-exporter"}
	app := NewApp(logger)

	exitCode := app.Run()
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestRunWithArgs(t *testing.T) {
	logger := newTestLogger(t)
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"vigilant-exporter", "some-arg"}

	app := NewApp(logger)

	exitCode := app.Run()
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func TestRunWithMultipleArgs(t *testing.T) {
	logger := newTestLogger(t)
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"vigilant-exporter", "arg1", "arg2", "arg3"}

	app := NewApp(logger)

	exitCode := app.Run()
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func newTestLogger(t *testing.T) *slog.Logger {
	t.Helper()

	handler := slog.NewJSONHandler(os.Stdout, nil)
	handler.WithAttrs([]slog.Attr{
		slog.String("test", t.Name()),
	})

	return slog.New(handler)
}
