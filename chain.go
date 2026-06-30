package httpmw

import "net/http"

// Middleware wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain composes middlewares from left to right.
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			if middlewares[i] != nil {
				next = middlewares[i](next)
			}
		}
		return next
	}
}

// Wrap applies a middleware chain to a handler.
func Wrap(next http.Handler, middlewares ...Middleware) http.Handler {
	return Chain(middlewares...)(next)
}
