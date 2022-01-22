package apiserver

import "net/http"

type ResponseWriter struct {
	// unonimous field - let us use do not release all methods in ResponseWriter; they will be accessable
	http.ResponseWriter
	// extention
	code int
}

// redefine internal method - WriteHeader
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode
	rw.ResponseWriter.WriteHeader(statusCode) // Continue with original method
}
