package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var service_name = viper.GetString("service.name")
var tracer = otel.Tracer(service_name)
var meter = otel.Meter(service_name)

func init() {
	apiCounter, err := meter.Int64Counter(
		"api.counter",
		metric.WithDescription("Number of API calls."),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		slog.Error("Failed to create metric counter `api.counter`.", "error", err)
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "api-counter-demo")
		defer span.End()

		span.AddEvent("Received HTTP request to '/'")
		slog.InfoContext(ctx, "Received HTTP request", "span-id", span.SpanContext().SpanID())

		apiCounter.Add(r.Context(), 1)
		w.Header().Add("x-demo", "12345")
		w.Write([]byte("Hello, world!"))
	})
}

func NewServer(ctx context.Context) {
	if viper.GetBool("server.tls.enabled") {
		slog.Debug("TLS enabled, loading certificates...")
	}

	slog.InfoContext(
		ctx,
		"Starting HTTP server.",
		"host", "localhost",
		"port", 8080,
		"tls", false,
	)
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		slog.Error("Failed to start HTTP server.", "error", err)
	}
}
