package web

import (
	"context"
	"fmt"
	"net/http"
)

type Handler func(context.Context, http.ResponseWriter, *http.Request) error

type Middleware func(Handler) Handler

type Logger func(context.Context, string, ...any)

type web struct {
	mux *http.ServeMux
}

func New() *web {
	return &web{
		mux: http.NewServeMux(),
	}
}

func (web *web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	web.mux.ServeHTTP(w, r)
}

func (web *web) Handle(method string, path string, handler Handler, middlewares ...Middleware) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		mid := middlewares[i]
		handler = mid(handler)
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		handler(r.Context(), w, r)
	}

	web.mux.HandleFunc(fmt.Sprintf("%s %s", method, path), h)
}
