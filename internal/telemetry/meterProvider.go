package telemetry

import (
	"context"
	"log/slog"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func NewMeterProvider(ctx context.Context) (metricProvider *metric.MeterProvider) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(viper.GetString("service.name")),
		),
	)
	if err != nil {
		slog.Error("Failed creating resource for metrics provider.", "error", err)
		return nil
	}

	metricExporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		slog.Error("Failed to initialize metrics exporter.", "error", err)
		return nil
	}

	metricProvider = metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
	)
	slog.Info("Successfully initialized the metrics provider.")
	return metricProvider
}
