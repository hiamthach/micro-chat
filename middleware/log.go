package middleware

import (
	"log"
	"net/http"
	"time"
)

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		defer func() {
			log.Printf(
				"%s %s %s %v",
				r.Method,
				r.RequestURI,
				r.RemoteAddr,
				time.Since(startTime),
			)
		}()

		h.ServeHTTP(w, r)
	})
}
