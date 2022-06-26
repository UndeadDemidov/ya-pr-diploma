package errors

import "errors"

var (
	ErrLoginIsInUseAlready = errors.New("given login is in use already")
)
