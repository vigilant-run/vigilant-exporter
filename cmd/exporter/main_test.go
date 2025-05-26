package main

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestRunWithNoArgs(t *testing.T) {
	logger := zap.NewNop()
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
	logger := zap.NewNop()
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
	logger := zap.NewNop()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"vigilant-exporter", "arg1", "arg2", "arg3"}

	app := NewApp(logger)

	exitCode := app.Run()
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}
