package main

import (
	"log"
	"net/http"
	"time"
)

type Middleware struct{}

func (m Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			log.Printf("[%s] %q %v", r.Method, r.URL.String(), time.Since(start))
		}(time.Now())
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (m Middleware) RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recover form panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
