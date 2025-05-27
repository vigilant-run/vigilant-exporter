package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ExporterConfig struct {
	FilePath string `mapstructure:"file"`
	Token    string `mapstructure:"token"`
	Endpoint string `mapstructure:"endpoint"`
	Insecure bool   `mapstructure:"insecure"`
}

func NewExporterConfig() *ExporterConfig {
	return &ExporterConfig{
		FilePath: "",
		Token:    "",
		Endpoint: "ingress.vigilant.run",
		Insecure: false,
	}
}

func InitConfig(cmd *cobra.Command) (*ExporterConfig, error) {
	viper.BindPFlag("file", cmd.Flags().Lookup("file"))
	viper.BindPFlag("token", cmd.Flags().Lookup("token"))
	viper.BindPFlag("endpoint", cmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("insecure", cmd.Flags().Lookup("insecure"))

	viper.SetDefault("endpoint", "ingress.vigilant.run")
	viper.SetDefault("insecure", false)

	config := &ExporterConfig{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

func AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("file", "f", "", "Path to the log file to monitor")
	cmd.Flags().StringP("token", "t", "", "Authentication token")
	cmd.Flags().StringP("endpoint", "e", "ingress.vigilant.run", "Endpoint URL for log ingestion")
	cmd.Flags().BoolP("insecure", "i", false, "Send logs over HTTP instead of HTTPS")
	cmd.Flags().SortFlags = false
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("token")
}
