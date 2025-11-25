package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var service_name = viper.GetString("service.name")
var tracer = otel.Tracer(service_name)
var meter = otel.Meter(service_name)

var apiCounter metric.Int64Counter

func init() {
	var err error
	apiCounter, err = meter.Int64Counter(
		"api.counter",
		metric.WithDescription("Amount of API calls"),
		metric.WithUnit("{call}"),
	)
	if err != nil {
		slog.Error("Failed initialize API Counter `api.counter` meter.", "error", err)
	}
}

func NewHTTPHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	return router
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
	defer span.End()

	span.AddEvent("Request started", trace.WithTimestamp(time.Now()))
	slog.InfoContext(
		ctx, "Received incoming request.",
		"path", r.URL.Path,
	)
	apiCounter.Add(ctx, 1)

	w.WriteHeader(http.StatusOK)
}
