package middlerwares

import (
	"net/http"
)

func Security_middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Pretection", "1:mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000;includeSubDomains;preload")
		w.Header().Set("Content-Security-Policy", "no-referrer")
		w.Header().Set("X-Powered-By", "Python/Django")

		next.ServeHTTP(w, r)
	})
}
