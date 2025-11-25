package telemetry

import (
	"context"
	"log/slog"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewLoggerProvider(ctx context.Context, mode TelemetryMode) (provider *log.LoggerProvider) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(viper.GetString("service.name")),
		),
	)
	if err != nil {
		slog.Error("Failed to initialize resource for logger provider.", "error", err)
		return nil
	}

	var exporter log.Exporter
	exporter, err = otlploggrpc.New(ctx)
	if mode == HTTP {
		exporter, err = otlploghttp.New(ctx)
	}
	if err != nil {
		slog.Error("Failed to initialize logger exporter", "error", err)
		return nil
	}

	processor := log.NewBatchProcessor(exporter)
	provider = log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(processor),
	)
	slog.Info("Successfully initialized the logger provider.")
	return provider
}
