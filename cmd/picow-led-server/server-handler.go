package main

import (
	"log/slog"
	"net/http"
)

type serverHandler struct{}

func (*serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered", "r", r)
		}
	}()

	crw := &serverResponseWriter{
		ResponseWriter: w,
		Hijacker:       w.(http.Hijacker),
	}

	crw.Header().Set("Access-Control-Allow-Origin", "*")
	http.DefaultServeMux.ServeHTTP(crw, r)

	log := slog.Warn

	if crw.status >= 200 && crw.status < 300 {
		if crw.Header().Get("Content-Type") == "" {
			crw.Header().Set("Content-Type", "application/json")
		}

		log = slog.Info
	} else if crw.status >= 500 || crw.status == 0 {
		log = slog.Error
	}

	log("Request",
		"status", crw.status,
		"addr", r.RemoteAddr,
		"method", r.Method,
		"url", r.URL,
	)
}
