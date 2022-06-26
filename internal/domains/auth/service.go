package auth

import (
	"context"

	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/handler"
)

type CredentialValidator interface {
	IsValid()
}

type CredentialManager interface {
	New(user, pword string) (cred CredentialValidator, ok bool)
	Add(cred CredentialValidator) error
	Get(user string) (cred CredentialValidator, err error)
	// Remove(credential CredentialValidator) error
}

var _ handler.Authenticator = (*Service)(nil)

type Service struct {
}

func (s Service) SignIn(ctx context.Context, user, pword string) error {
	return errors2.ErrLoginIsInUseAlready
}

func (s Service) Login(ctx context.Context, user, pword string) error {
	// TODO implement me
	panic("implement me")
}
