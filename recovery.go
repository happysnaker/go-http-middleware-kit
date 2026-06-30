package httpmw

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

type RecoveryOptions struct {
	Logger            *slog.Logger
	DisableStackTrace bool
	Handler           func(http.ResponseWriter, *http.Request, any)
}

// Recovery converts panics into HTTP 500 responses and logs the panic details.
func Recovery(opts RecoveryOptions) Middleware {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	handler := opts.Handler
	if handler == nil {
		handler = func(w http.ResponseWriter, _ *http.Request, _ any) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}

				fields := []any{
					"panic", recovered,
					"request_id", GetRequestID(r.Context()),
					"client_ip", GetClientIP(r.Context()),
					"method", r.Method,
					"path", r.URL.Path,
				}
				if !opts.DisableStackTrace {
					fields = append(fields, "stack", string(debug.Stack()))
				}

				logger.Error("http request panicked", fields...)
				handler(w, r, recovered)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
