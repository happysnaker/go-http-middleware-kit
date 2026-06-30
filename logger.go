package httpmw

import (
	"log/slog"
	"net/http"
	"time"
)

type RequestLoggerOptions struct {
	Logger           *slog.Logger
	IncludeUserAgent bool
}

// RequestLogger writes a structured log entry for each completed request.
func RequestLogger(opts RequestLoggerOptions) Middleware {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			recorder := newResponseRecorder(w)

			next.ServeHTTP(recorder, r)

			fields := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", recorder.statusCode,
				"bytes", recorder.bytesWritten,
				"latency", time.Since(startedAt).String(),
				"request_id", GetRequestID(r.Context()),
				"client_ip", GetClientIP(r.Context()),
			}
			if opts.IncludeUserAgent {
				fields = append(fields, "user_agent", r.UserAgent())
			}

			logger.Info("http request completed", fields...)
		})
	}
}
