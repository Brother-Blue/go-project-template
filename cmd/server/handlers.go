package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

func tracerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()
		// Must set the request context to the context that contains the span, this will be passed down through proceeding middlewares
		r = r.WithContext(ctx)
		span.SetStatus(codes.Ok, "")
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(
			r.Context(), "Incoming request",
			"http.path", r.URL.Path,
			"http.method", r.Method,
			"http.scheme", r.URL.Scheme,
			"user_agent", r.UserAgent(),
			"headers", r.Header,
			"cookies", r.Cookies(),
		)
		next.ServeHTTP(w, r)
	})
}

func NewHTTPHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)

	router.Use(tracerMiddleware)
	router.Use(loggingMiddleware)
	return router
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Starting a new trace that will receive the parent trace from the tracerMiddleware via r.Context()
	ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
	defer span.End()

	if r.Method != http.MethodGet {
		span.SetAttributes(attribute.String("http.method", r.Method))
		slog.WarnContext(ctx, "Received bad request.", "http.method", r.Method)
		span.RecordError(http.ErrBodyNotAllowed, trace.WithStackTrace(true))
		span.SetStatus(codes.Error, "Unaccepted method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Not ok"))
		return
	}

	span.AddEvent("Request started", trace.WithTimestamp(time.Now()))
	apiCounter.Add(ctx, 1)
	span.SetStatus(codes.Ok, "")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}
