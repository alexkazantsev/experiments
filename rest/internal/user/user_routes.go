package user

import (
	"github.com/alexkazantsev/experiments/rest/server"
)

func Routes(router *server.Router, ctrl UserController) {
	router.V1.HandleFunc("GET /users", ctrl.Find)
	router.V1.HandleFunc("GET /users/{id}", ctrl.FindOne)
}
