package httpmw

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DefaultRequestIDHeader = "X-Request-ID"

type requestIDContextKey struct{}

type RequestIDOptions struct {
	HeaderName    string
	Generator     func() string
	TrustIncoming bool
}

func (o RequestIDOptions) normalized() RequestIDOptions {
	if o.HeaderName == "" {
		o.HeaderName = DefaultRequestIDHeader
	}
	if o.Generator == nil {
		o.Generator = NewRequestID
	}
	return o
}

// RequestID injects a request id into the request context and response header.
func RequestID(opts RequestIDOptions) Middleware {
	opts = opts.normalized()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := ""
			if opts.TrustIncoming {
				id = strings.TrimSpace(r.Header.Get(opts.HeaderName))
			}
			if id == "" {
				id = opts.Generator()
			}

			w.Header().Set(opts.HeaderName, id)
			ctx := context.WithValue(r.Context(), requestIDContextKey{}, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// NewRequestID generates a compact random request id.
func NewRequestID() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(buf[:])
}

// GetRequestID returns the request id stored in the context.
func GetRequestID(ctx context.Context) string {
	value, _ := ctx.Value(requestIDContextKey{}).(string)
	return value
}
