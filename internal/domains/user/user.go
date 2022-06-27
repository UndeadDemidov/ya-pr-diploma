package user

import "github.com/google/uuid"

type User struct {
	ID string
}

func NewUser() User {
	return User{ID: uuid.New().String()}
}
