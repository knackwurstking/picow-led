package main

import "net/http"

type customResponseWriter struct {
	http.ResponseWriter
	http.Hijacker

	status int
}

func (crw *customResponseWriter) WriteHeader(statusCode int) {
	crw.status = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}
