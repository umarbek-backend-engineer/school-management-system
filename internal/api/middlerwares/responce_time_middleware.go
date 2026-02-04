package middlerwares

import (
	"log"
	"net/http"
	"time"
)

func Responce_time(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// creating custom responce writer to caputer the status code
		wrapped_writer := &response{ResponseWriter: w, status: http.StatusOK}

		duration := time.Since(start)
		w.Header().Set("X-Responce-Time", duration.String())
		next.ServeHTTP(wrapped_writer, r)
		duration = time.Since(start)
		log.Printf("Method: %s, URL: %s, Status:%d, Duration: %v\n", r.Method, r.URL, wrapped_writer.status, duration.String())
	})
}

type response struct {
	http.ResponseWriter
	status int
}

func (rw *response) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
