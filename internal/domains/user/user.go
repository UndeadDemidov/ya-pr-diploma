package user

import "github.com/google/uuid"

type User struct {
	ID string
}

func NewUser() User {
	// ToDo сохранить в репозитории!
	return User{ID: uuid.New().String()}
}
