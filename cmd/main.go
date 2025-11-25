package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/brother-blue/campaign-manager/cmd/server"
	_ "github.com/brother-blue/campaign-manager/internal/config"
	"github.com/brother-blue/campaign-manager/internal/telemetry"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
)

func main() {
	ctx := context.Background()

	if viper.GetBool("config.telemetry.enabled") {
		tracerProvider := telemetry.NewTracerProvider(ctx)
		meterProvider := telemetry.NewMeterProvider(ctx)
		loggerProvider := telemetry.NewLoggerProvider(ctx)
		global.SetLoggerProvider(loggerProvider)

		if tracerProvider == nil || meterProvider == nil || loggerProvider == nil {
			slog.Error("Failed to initialize all telemetry providers. Please see previous errors for potential causes.")
			os.Exit(1)
		}

		logger := otelslog.NewLogger(viper.GetString("service.name") + "-logger")

		defer tracerProvider.Shutdown(ctx)
		defer meterProvider.Shutdown(ctx)
		defer loggerProvider.Shutdown(ctx)

		otel.SetTracerProvider(tracerProvider)
		otel.SetMeterProvider(meterProvider)
		slog.SetDefault(logger)
	}

	if err := server.Start(ctx); err != nil {
		slog.Error("Error received from server.", "error", err)
	}
}
