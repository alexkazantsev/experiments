package middlewares

import "net/http"

func ContentType(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(writer, request)
	}
}
