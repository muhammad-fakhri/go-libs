package httpmiddleware

import (
	"bytes"
	"net/http"
)

type LogResponseWriter struct {
	statusCode int
	bodyBuffer *bytes.Buffer
	http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{
		ResponseWriter: w,
		// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
		// we default to that status code.
		statusCode: http.StatusOK,
		bodyBuffer: &bytes.Buffer{},
	}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.bodyBuffer.Write(body)
	return w.ResponseWriter.Write(body)
}
