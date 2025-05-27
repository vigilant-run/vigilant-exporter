package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"vigilant-exporter/internal/app"
	"vigilant-exporter/internal/config"

	"github.com/spf13/cobra"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rootCmd := newRootCmd(ctx)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(ctx context.Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vigilant-exporter",
		Short: "A log file exporter for Vigilant",
		Long:  "Monitors log files and exports them to a Vigilant endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := config.InitConfig(cmd)
			if err != nil {
				return fmt.Errorf("configuration error: %w", err)
			}

			app := app.NewApp(config)
			return app.Run(ctx)
		},
	}

	config.AddFlags(rootCmd)

	return rootCmd
}
