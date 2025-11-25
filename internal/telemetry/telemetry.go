package telemetry

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TelemetryMode string

const (
	GRPC TelemetryMode = "grpc"
	HTTP TelemetryMode = "http"
)

type TelemetryProvider struct {
	LogProvider    *log.LoggerProvider
	MeterProvider  *metric.MeterProvider
	TracerProvider *trace.TracerProvider
}

func NewTelemetryProvider(ctx context.Context, mode TelemetryMode) TelemetryProvider {
	loggerProvider := NewLoggerProvider(ctx, mode)
	meterProvider := NewMeterProvider(ctx, mode)
	tracerProvider := NewTracerProvider(ctx, mode)

	if tracerProvider == nil || meterProvider == nil || loggerProvider == nil {
		slog.Error("Failed to initialize all telemetry providers. Please see previous errors for potential causes.")
		os.Exit(1)
	}

	global.SetLoggerProvider(loggerProvider)
	logger := otelslog.NewLogger(viper.GetString("service.name") + "-logger")
	slog.SetDefault(logger)
	otel.SetMeterProvider(meterProvider)
	otel.SetTracerProvider(tracerProvider)

	return TelemetryProvider{
		LogProvider:    loggerProvider,
		MeterProvider:  meterProvider,
		TracerProvider: tracerProvider,
	}
}

func (t *TelemetryProvider) Stop(ctx context.Context) {
	t.LogProvider.Shutdown(ctx)
	t.MeterProvider.Shutdown(ctx)
	t.TracerProvider.Shutdown(ctx)
}
