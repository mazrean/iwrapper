package example

import (
	"bufio"
	"net"
	"net/http"
)

type MyResponseWriter struct {
	http.ResponseWriter
	size *int
}

func WrapResponseWriter(w http.ResponseWriter) (http.ResponseWriter, *int) {
	size := 0
	return ResponseWriterWrapper(w, func(w http.ResponseWriter) ResponseWriter {
		return MyResponseWriter{w, &size}
	}), &size
}

func (w MyResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	*w.size += n
	return n, err
}

func (w MyResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w MyResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w MyResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}
