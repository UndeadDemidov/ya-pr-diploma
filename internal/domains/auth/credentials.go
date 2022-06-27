package auth

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
)

type Repository interface {
	Create(ctx context.Context, user user.User, login, pword string) error
	Read(ctx context.Context, login string) (user user.User, err error)
}

var _ CredentialManager = (*Credentials)(nil)

// Никакого хранения в памяти кредов не делаем. Чем меньше мест, где храняться логины/пароли,
// тем меньше мест нужно защищать от хакеров - даже если это дополнительное место - память.
// Имеет смысл кешировать столько с очень большой нагрузкой по аутентификации.
type Credentials struct {
	repo Repository
}

func (c *Credentials) AddNewUser(ctx context.Context, usr user.User, login, pword string) error {
	// ToDo сделать тут хеш пароля с солью
	return c.repo.Create(ctx, usr, login, pword)
}

func (c *Credentials) GetUser(ctx context.Context, login string) (usr user.User, err error) {
	// ToDo сделать тут хеш пароля с солью
	// Нужно сделать отдельно Get и отдельно Authenticate!
	usr, err = c.repo.Read(ctx, login)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (c *Credentials) AuthenticateUser(ctx context.Context, login, pword string) (usr user.User, err error) {
	// TODO implement me
	panic("implement me")
}
