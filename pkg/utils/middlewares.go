package utils

import "net/http"

type Middlewares func(http.Handler) http.Handler

func ApplayMiddlewares(handler http.Handler, mid ...Middlewares) http.Handler {
	for _, v := range mid {
		handler = v(handler)
	}
	return handler
}
