package middlerwares

import (
	"net/http"
)

var allowed_origins = []string{
	"http://localhost:8080",
	"http://localhost:3000",
	"https://localhost:8080",
	"https://localhost:3000",
	"null", // <-- VERY IMPORTANT for file://
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			http.Error(w, "Not allowed by cors middleware", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})

}

func isOriginAllowed(origin string) bool {
	for _, v := range allowed_origins {
		if origin == v {
			return true
		}
	}
	return false
}
