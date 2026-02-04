package middlerwares

import (
	"net/http"
	"strings"
)

func Exclude_Routes(middleware func(http.Handler) http.Handler, excludedPath ...string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range excludedPath {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}
			middleware(next).ServeHTTP(w, r)
		})
	}

}
