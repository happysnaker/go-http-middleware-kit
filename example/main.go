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
		_, _ = w.Write([]byte("hello, middleware kit"))
	})

	handler := httpmw.Wrap(
		mux,
		httpmw.RealIP(httpmw.RealIPOptions{SetRemoteAddr: true}),
		httpmw.RequestID(httpmw.RequestIDOptions{TrustIncoming: true}),
		httpmw.Recovery(httpmw.RecoveryOptions{Logger: logger}),
		httpmw.RequestLogger(httpmw.RequestLoggerOptions{Logger: logger, IncludeUserAgent: true}),
		httpmw.Timeout(httpmw.TimeoutOptions{Duration: 5 * time.Second}),
	)

	logger.Info("server starting", "addr", ":8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.Error("server stopped", "error", err)
	}
}
