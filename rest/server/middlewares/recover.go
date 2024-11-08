package middlewares

import (
	"errors"
	"log"
	"net/http"
)

func Recover(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var err error

				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("undefined error")
				}

				log.Printf("recover from panic: %+v", err)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(writer, request)
	}
}
