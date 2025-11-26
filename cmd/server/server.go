package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/viper"
)

func Start(ctx context.Context) (err error) {
	if viper.GetBool("server.tls.enabled") {
		slog.InfoContext(ctx, "TLS enabled, loading certificates...")
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	server := &http.Server{
		Addr:    ":8080",
		Handler: NewHTTPHandler(),
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Starting HTTP server on :8080")
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err = <-serverErrors:
		return err
	case <-ctx.Done():
		stop()
	}
	err = server.Shutdown(ctx)
	return err
}
