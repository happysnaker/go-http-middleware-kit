package httpmw

import (
	"context"
	"net"
	"net/http"
	"strings"
)

type clientIPContextKey struct{}

type RealIPOptions struct {
	HeaderNames   []string
	SetRemoteAddr bool
}

func (o RealIPOptions) normalized() RealIPOptions {
	if len(o.HeaderNames) == 0 {
		o.HeaderNames = []string{"X-Forwarded-For", "X-Real-IP"}
	}
	return o
}

// RealIP extracts the client IP from trusted proxy headers and stores it in context.
func RealIP(opts RealIPOptions) Middleware {
	opts = opts.normalized()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractClientIP(r, opts.HeaderNames)
			if ip == "" {
				ip = normalizeRemoteAddr(r.RemoteAddr)
			}

			ctx := context.WithValue(r.Context(), clientIPContextKey{}, ip)
			updated := r.WithContext(ctx)
			if opts.SetRemoteAddr && ip != "" {
				clone := updated.Clone(ctx)
				clone.RemoteAddr = net.JoinHostPort(ip, "0")
				updated = clone
			}

			next.ServeHTTP(w, updated)
		})
	}
}

// GetClientIP returns the client IP stored in the context.
func GetClientIP(ctx context.Context) string {
	value, _ := ctx.Value(clientIPContextKey{}).(string)
	return value
}

func extractClientIP(r *http.Request, headers []string) string {
	for _, name := range headers {
		raw := strings.TrimSpace(r.Header.Get(name))
		if raw == "" {
			continue
		}

		if strings.EqualFold(name, "X-Forwarded-For") {
			for _, part := range strings.Split(raw, ",") {
				if ip := normalizeIP(part); ip != "" {
					return ip
				}
			}
			continue
		}

		if ip := normalizeIP(raw); ip != "" {
			return ip
		}
	}

	return ""
}

func normalizeRemoteAddr(remoteAddr string) string {
	remoteAddr = strings.TrimSpace(remoteAddr)
	if remoteAddr == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return host
	}
	return normalizeIP(remoteAddr)
}

func normalizeIP(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(raw); err == nil {
		raw = host
	}

	ip := net.ParseIP(strings.Trim(raw, "[]"))
	if ip == nil {
		return ""
	}
	return ip.String()
}
