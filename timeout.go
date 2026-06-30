package httpmw

import (
	"net/http"
	"time"
)

type TimeoutOptions struct {
	Duration time.Duration
	Message  string
}

// Timeout wraps the handler with http.TimeoutHandler.
func Timeout(opts TimeoutOptions) Middleware {
	if opts.Duration <= 0 {
		opts.Duration = 30 * time.Second
	}
	if opts.Message == "" {
		opts.Message = http.StatusText(http.StatusGatewayTimeout)
	}

	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, opts.Duration, opts.Message)
	}
}
