package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/nghialm269/fault-demo/pkg/fserrors/wrappers/errctx"
	"github.com/nghialm269/fault-demo/pkg/fserrors/wrappers/errstacktrace"
	"github.com/nghialm269/fault-demo/pkg/web"
)

func Logger(logger *slog.Logger) web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next nextFunc) error {
		now := time.Now()

		method := r.Method
		path := r.URL.Path
		remoteAddr := r.RemoteAddr

		logger = logger.With("method", method, "path", path, "remoteaddr", remoteAddr)

		logger.InfoContext(ctx, "request started")

		err := next(ctx)

		logger = logger.With("since", time.Since(now).String())

		if err != nil {
			logger.ErrorContext(ctx, "request completed", "error", err, "errctx", errctx.Unwrap(err), "stacktrace", errstacktrace.Get(err))
			return err
		}

		logger.InfoContext(ctx, "request completed")

		return nil
	}

	return newMiddleware(midFunc)
}
