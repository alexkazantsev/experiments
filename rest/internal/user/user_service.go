package user

import (
	"context"
	"log"

	"github.com/alexkazantsev/experiments/rest/domain"
	"github.com/google/uuid"
)

type UserService interface {
	FindOne(ctx context.Context, id uuid.UUID) *domain.User
}

type UserServiceImpl struct {
}

func (u UserServiceImpl) FindOne(ctx context.Context, id uuid.UUID) *domain.User {
	log.Printf("user id: %s\n", ctx.Value("user"))

	return &domain.User{
		ID:    uuid.New(),
		Name:  "Foo",
		Email: "foo@bar.com",
	}
}

func NewUserService() UserService {
	return &UserServiceImpl{}
}
