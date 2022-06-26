package auth

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

var _ CredentialManager = (*Credentials)(nil)

// Никакого хранения в памяти кредов не делаем. Чем меньше мест, где храняться логины/пароли,
// тем меньше мест нужно защищать от хакеров - даже если это дополнительное место - память.
// Имеет смысл кешировать столько с очень большой нагрузкой по аутентификации.
type Credentials struct {
}

func (c *Credentials) AddNewUser(ctx context.Context, user user.User, login, pword string) error {
	// TODO implement me
	// ToDo сделать тут хеш пароля с солью
	// создаем пользователя, сохраняем его в БД
	// добавляем ему креды
	panic("implement me")
}

func (c *Credentials) GetUser(ctx context.Context, login, pword string) (user user.User, err error) {
	// ToDo сделать тут хеш пароля с солью
	// TODO implement me
	panic("implement me")
}
