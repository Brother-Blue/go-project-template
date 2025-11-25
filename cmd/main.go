package main

import (
	"context"
	"log/slog"

	"github.com/brother-blue/campaign-manager/cmd/server"
	_ "github.com/brother-blue/campaign-manager/internal/config"
	"github.com/brother-blue/campaign-manager/internal/telemetry"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()

	if viper.GetBool("config.telemetry.enabled") {
		tp := telemetry.NewTelemetryProvider(
			ctx,
			telemetry.TelemetryMode(viper.GetString("config.telemetry.mode")),
		)
		defer tp.Stop(ctx)
	}

	if err := server.Start(ctx); err != nil {
		slog.Error("Error received from server.", "error", err)
	}
}
