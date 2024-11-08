package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var t = time.Now()

		log.Printf("-> [REQ] Method: %s, Path: %s, Time Start: %s",
			request.Method, request.URL.Path, t.Format(time.DateTime))

		next.ServeHTTP(writer, request)

		log.Printf("<- [RES] Method: %s, Path: %s, Time End: %s, Tooks: %dµs",
			request.Method, request.URL.Path, t.Format(time.DateTime), time.Since(t).Microseconds())
	}
}
