package response

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

const (
	statusCodeUninitialized = -1
)

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     statusCodeUninitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) GetStatusCode() int {
	return rw.statusCode
}

func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == statusCodeUninitialized {
		panic("no status code set")
	}

	return rw.statusCode
}
