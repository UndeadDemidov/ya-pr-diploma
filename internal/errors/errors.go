package errors

import "github.com/pkg/errors"

var (
	ErrLoginIsInUseAlready      = errors.New("given login is in use already")
	ErrPairLoginPwordIsNotExist = errors.New("given pair login and password is not exists")
	ErrSessionIsExpired         = errors.New("session is expired")
)
