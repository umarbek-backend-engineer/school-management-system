package middlerwares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		w = &gzipresponse{ResponseWriter: w, Writer: gz}

		next.ServeHTTP(w, r)
	})
}

type gzipresponse struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (gr *gzipresponse) Write(b []byte) (int, error) {
	return gr.Writer.Write(b)
}
