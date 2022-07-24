package errors

import "github.com/pkg/errors"

// SignIn/Login errors
var (
	ErrLoginIsInUseAlready      = errors.New("given login is in use already")
	ErrPairLoginPwordIsNotExist = errors.New("given pair login and password is not exists")
)

// Sessions errors
var (
	ErrSessionIsExpired           = errors.New("session is expired")
	ErrSessionUserCanNotBeDefined = errors.New("user can't be defined by session")
)

// Orders errors
var (
	ErrOrderAlreadyUploaded              = errors.New("order is already uploaded by this user")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order is already uploaded by another user")
	ErrOrderInvalidNumberFormat          = errors.New("invalid order number format")
)
