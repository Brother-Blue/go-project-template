package telemetry

import (
	"context"
	"log/slog"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewTracerProvider(ctx context.Context) (provider *sdktrace.TracerProvider) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(viper.GetString("service.name")),
		),
	)
	if err != nil {
		slog.Error("Failed creating resources for trace provider.", "error", err)
		return nil
	}

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		slog.Error("Failed to create OTLP exporter.")
		return nil
	}

	provider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	slog.Info("Successfully initialized the trace provider.")
	return provider
}
