package auth

import (
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/handler"
)

var _ handler.CredentialManager = (*Credentials)(nil)

type Credentials struct {
}

func (c Credentials) New(user, pword string) (cred handler.CredentialValidator, ok bool) {
	// TODO implement me
	panic("implement me")
}

func (c Credentials) Add(cred handler.CredentialValidator) error {
	// TODO implement me
	panic("implement me")
}

func (c Credentials) Get(user string) (cred handler.CredentialValidator, err error) {
	// TODO implement me
	panic("implement me")
}
