package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nghialm269/fault-demo/pkg/fserrors"
	"github.com/nghialm269/fault-demo/pkg/web"
)

func Panics() web.Middleware {
	midFunc := func(ctx context.Context, r *http.Request, next nextFunc) (err error) {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("panic: %+v", r)
				err = fserrors.New(message)

			}
		}()

		return next(ctx)
	}

	return newMiddleware(midFunc)
}
