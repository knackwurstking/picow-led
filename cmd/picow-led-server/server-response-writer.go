package main

import "net/http"

type serverResponseWriter struct {
	http.ResponseWriter
	http.Hijacker

	status int
}

func (crw *serverResponseWriter) WriteHeader(statusCode int) {
	crw.status = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}
