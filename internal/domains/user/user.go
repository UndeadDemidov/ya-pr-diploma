package user

import (
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/google/uuid"
)

var _ middleware.Referencer = (*User)(nil)

type User struct {
	ID string
}

func NewUser() User {
	return User{ID: uuid.New().String()}
}

func (u User) Reference() string {
	return u.ID
}
