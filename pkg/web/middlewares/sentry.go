package middlewares

import (
	"context"
	"net/http"

	"github.com/getsentry/sentry-go"

	"github.com/nghialm269/fault-demo/pkg/fserrors/wrappers/errctx"
	"github.com/nghialm269/fault-demo/pkg/web"
)

func Sentry() web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next nextFunc) error {
		err := next(ctx)

		if err != nil {
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetContext("errctx", errctx.Unwrap(err))
				scope.SetLevel(sentry.LevelError)

				sentry.CaptureException(err)
			})
		}

		return err

	}

	return newMiddleware(midFunc)
}
