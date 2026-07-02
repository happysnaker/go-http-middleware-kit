# go-http-middleware-kit

[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/github/license/happysnaker/go-http-middleware-kit)](./LICENSE)
[![Stars](https://img.shields.io/github/stars/happysnaker/go-http-middleware-kit?style=social)](https://github.com/happysnaker/go-http-middleware-kit/stargazers)
[![Support](https://img.shields.io/badge/support-WeChat%20%26%20Alipay-7aa2ff)](https://happysnaker.github.io/support/#from-go-http-middleware-kit)
[![Async Review](https://img.shields.io/badge/review-Quick%20read%20%2F%20async-9b87f5)](https://happysnaker.github.io/review/)

Reusable `net/http` middleware for backend services that want **request IDs**, **real IP extraction**, **structured logs**, **panic recovery**, and **timeouts** without pulling in a full framework.

This repo is intentionally small, readable, and production-minded:

- standard-library first
- easy to copy into internal services
- composable around ordinary `http.Handler`
- useful for starter repos, internal APIs, and interview-friendly service demos

- Project page: [happysnaker.github.io/go-http-middleware-kit](https://happysnaker.github.io/go-http-middleware-kit/)

## What is included

- `RequestID` middleware with context + response-header propagation
- `RealIP` middleware for `X-Forwarded-For` / `X-Real-IP`
- `RequestLogger` middleware built on `log/slog`
- `Recovery` middleware with panic logging and stack traces
- `Timeout` middleware backed by `http.TimeoutHandler`
- `Chain` / `Wrap` helpers for clean composition

## Quick start

```bash
go get github.com/happysnaker/go-http-middleware-kit
```

```go
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	httpmw "github.com/happysnaker/go-http-middleware-kit"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	})

	handler := httpmw.Wrap(
		mux,
		httpmw.RealIP(httpmw.RealIPOptions{SetRemoteAddr: true}),
		httpmw.RequestID(httpmw.RequestIDOptions{TrustIncoming: true}),
		httpmw.Recovery(httpmw.RecoveryOptions{Logger: logger}),
		httpmw.RequestLogger(httpmw.RequestLoggerOptions{Logger: logger, IncludeUserAgent: true}),
		httpmw.Timeout(httpmw.TimeoutOptions{Duration: 5 * time.Second}),
	)

	_ = http.ListenAndServe(":8080", handler)
}
```

## Middleware overview

| Middleware | Purpose |
| --- | --- |
| `RequestID` | attach / generate request ids and expose them through context |
| `RealIP` | resolve client IP from trusted proxy headers |
| `RequestLogger` | emit structured request completion logs via `slog` |
| `Recovery` | catch panics, log stack traces, and return HTTP 500 |
| `Timeout` | cap slow handlers with a timeout response |

## Why this repo exists

Many backend teams want the same small layer of HTTP middleware, but:

- a lot of examples are too toy-like to reuse
- some libraries pull in much more framework than needed
- starter repos often hide the middleware details instead of making them easy to read

`go-http-middleware-kit` aims for the middle ground: enough to be practical, small enough to understand in one sitting.

## Project layout

```text
chain.go            middleware composition helpers
request_id.go       request-id propagation and generation
real_ip.go          trusted proxy IP extraction
logger.go           slog request logging
recovery.go         panic recovery middleware
timeout.go          request timeout wrapper
response_writer.go  response recorder for logs / metrics-like data
example/main.go     minimal integration example
```

## Good fit for

- small internal APIs
- service starter templates
- backend interview demos
- teams that want lightweight `net/http` middleware without a framework migration

## Related repos

- [`go-service-starter`](https://github.com/happysnaker/go-service-starter) — minimal production-minded Go HTTP service starter
- [`backend-engineer-checklist`](https://github.com/happysnaker/backend-engineer-checklist) — practical backend growth roadmap
- [`system-design-checklist`](https://github.com/happysnaker/system-design-checklist) — reusable system-design checklist

## Support

If this repo saves you time, consider:

- starring the repo
- sharing it with other backend engineers
- supporting my open-source work via the [support page](https://happysnaker.github.io/support/#from-go-http-middleware-kit)
- if you want lightweight async feedback on a public GitHub profile, repo README, or portfolio page, details are also available there

If this middleware kit saved you copy-paste time in a real service, small direct support is especially appreciated:

- **¥9.9** — if one middleware or example saved you a quick detour
- **¥19.9** — if it helped you wire request IDs / logging / recovery faster
- **best payment note** — `go-http-middleware-kit` or `request-id middleware`
- **fastest path** — tip directly on the support page if this repo saved you time; use **¥29.9** / **¥99** only if you also want feedback back
- **¥99** — if you want compact async feedback on your own backend repo or README

**Fastest path if this repo helped:** open the [support page](https://happysnaker.github.io/support/#from-go-http-middleware-kit), scan WeChat / Alipay, and leave a short note like `go-http-middleware-kit` or `request-id middleware`. Tying the tip to one concrete repo tends to convert much better than a generic donation.

## License

MIT
