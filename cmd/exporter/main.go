package main

import (
	"fmt"
	"os"
	"vigilant-exporter/internal/app"
	"vigilant-exporter/internal/config"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
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
			return app.Run()
		},
	}

	config.AddFlags(rootCmd)

	return rootCmd
}
