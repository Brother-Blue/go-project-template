package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Startup logger
func init() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(l)
}

// Setup defaults
func init() {
	viper.SetDefault("otel_protocol", "http/protobuf")
	viper.SetDefault("otel_endpoint", "http://localhost:4318")
	viper.SetDefault("otel_insecure", true)
}

// Config setup and parsing
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$PWD")
	viper.AddConfigPath("$HOME/.campaign_manager")
	viper.AddConfigPath("/etc/campaign_manager")

	// Config parsing
	if err := viper.ReadInConfig(); err != nil {
		var configNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configNotFoundError) {
			slog.Error(
				"Config file could not be found at any of the specified paths.",
				"dev_path", "`.` - Only when the `developer_mode` config is enabled.",
				"home_path", "$HOME/.campaign_manager - The default path used.",
				"etc_path", "/etc/campaign_manager - The alternative path used (least precedence).",
			)
			os.Exit(1)
		}
		slog.Error(
			"Failed to read in config file.",
			"error", err,
		)
		os.Exit(1)
	}
	slog.Info(
		"Loaded config file.",
		"path", viper.GetViper().ConfigFileUsed(),
	)

	if viper.GetBool("developer_mode.enabled") {
		slog.Info("Developer mode enabled, using project's `config.yaml` and attempting to read in a .env file.")
		if err := godotenv.Load(); err != nil {
			slog.Warn("Failed loading in .env files. If you did not include any you can safely ignore this warning.")
		} else {
			slog.Info("Successfully read in .env file.")
		}
	}

	// Telemetry
	if viper.GetBool("config.telemetry.enabled") {
		os.Setenv("OTEL_SERVICE_NAME", viper.GetString("service.name"))
		os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", viper.GetString("config.telemetry.otlp_protocol"))
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", viper.GetString("config.telemetry.otlp_endpoint"))
		os.Setenv("OTEL_EXPORTER_OTLP_INSECURE", viper.GetString("config.telemetry.otlp_insecure"))

		slog.Info(
			"Telemetry has been configured.",
			"service_name", viper.GetString("service.name"),
			"otlp_endpoint", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
			"otlp_protocol", os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"),
			"provider_protocol", viper.GetString("config.telemetry.mode"),
		)
	}

	// Config hot reloading
	if viper.GetBool("config.hot_reload") {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			slog.Warn(
				"Config file updated. Reloading config, please see below for any errors. You may receive multiple reload messages, deduplcated events have not yet been implemented and I haven't gotten around to doing it yet.",
				"issue_url", "https://github.com/spf13/viper/issues/948",
			)
			if err := viper.ReadInConfig(); err != nil {
				slog.Error(
					"Failed to read in updated changes.",
					"error", err,
				)
			}
		})
	}
}
