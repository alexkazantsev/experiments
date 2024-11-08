package middlewares

import (
	"net/http"
	"time"
)

func Timeout(d time.Duration) func(http.Handler) http.HandlerFunc {
	return func(handler http.Handler) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			handler = http.TimeoutHandler(handler, d, "timout message")
			handler.ServeHTTP(writer, request)
		}
	}
}
