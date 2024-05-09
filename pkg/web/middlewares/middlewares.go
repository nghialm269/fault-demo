package middlewares

import (
	"context"
	"net/http"

	"github.com/nghialm269/fault-demo/pkg/web"
)

type nextFunc func(ctx context.Context) error
type midFunc func(context.Context, *http.Request, nextFunc) error

func newMiddleware(midFunc midFunc) web.Middleware {
	m := func(webHandler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			next := func(ctx context.Context) error {
				return webHandler(ctx, w, r)
			}

			return midFunc(ctx, r, next)
		}

		return h
	}

	return m
}
