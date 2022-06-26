package auth

var _ CredentialManager = (*Credentials)(nil)

type Credentials struct {
}

func (c Credentials) New(user, pword string) (cred CredentialValidator, ok bool) {
	// TODO implement me
	panic("implement me")
}

func (c Credentials) Add(cred CredentialValidator) error {
	// TODO implement me
	panic("implement me")
}

func (c Credentials) Get(user string) (cred CredentialValidator, err error) {
	// TODO implement me
	panic("implement me")
}
