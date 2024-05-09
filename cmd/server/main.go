package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/nghialm269/fault-demo/internal/users"
	"github.com/nghialm269/fault-demo/pkg/web"
	"github.com/nghialm269/fault-demo/pkg/web/middlewares"

	"github.com/getsentry/sentry-go"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	initSentry()
	defer sentry.Flush(2 * time.Second)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	web := web.New()
	web.Handle(
		http.MethodGet,
		"/users/{id}",
		users.HandlerGetUserByID,
		middlewares.Logger(logger),
		middlewares.Sentry(),
		middlewares.Panics(),
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: web,
	}

	logger.InfoContext(ctx, "startup", "host", server.Addr)
	server.ListenAndServe()
}

func initSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		panic(err)
	}
}
