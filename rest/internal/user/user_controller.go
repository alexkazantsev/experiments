package user

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type UserController interface {
	FindOne(http.ResponseWriter, *http.Request)
	Find(http.ResponseWriter, *http.Request)
}

type UserControllerImpl struct {
	service UserService
}

func (u UserControllerImpl) Find(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write([]byte("bar"))
}

func (u UserControllerImpl) FindOne(w http.ResponseWriter, r *http.Request) {
	var user = u.service.FindOne(r.Context(), uuid.MustParse(r.PathValue("id")))
	_ = json.NewEncoder(w).Encode(user)
}

func NewUserController(service UserService) UserController {
	return &UserControllerImpl{service: service}
}
