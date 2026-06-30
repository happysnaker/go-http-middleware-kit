package httpmw

import (
	"bufio"
	"errors"
	"io"
	"net"
	"net/http"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
	wroteHeader  bool
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}
	r.statusCode = statusCode
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(body []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(body)
	r.bytesWritten += n
	return n, err
}

func (r *responseRecorder) Flush() {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (r *responseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("httpmw: response writer does not support hijacking")
	}
	return hijacker.Hijack()
}

func (r *responseRecorder) Push(target string, opts *http.PushOptions) error {
	pusher, ok := r.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}
	return pusher.Push(target, opts)
}

func (r *responseRecorder) ReadFrom(src io.Reader) (int64, error) {
	readerFrom, ok := r.ResponseWriter.(io.ReaderFrom)
	if !ok {
		return io.Copy(r, src)
	}
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := readerFrom.ReadFrom(src)
	r.bytesWritten += int(n)
	return n, err
}
