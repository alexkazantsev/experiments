package user

import (
	"net/http"
)

func RegisterRoutes(v1 *http.ServeMux, ctrl UserController) {
	v1.HandleFunc("GET /users", ctrl.Find)
	v1.HandleFunc("GET /users/{id}", ctrl.FindOne)
}
