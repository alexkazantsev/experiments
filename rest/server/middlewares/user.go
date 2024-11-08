package middlewares

import (
	"context"
	"net/http"
)

func User(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var id = request.Header.Get("x-user")

		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), "user", id)))
	}
}
