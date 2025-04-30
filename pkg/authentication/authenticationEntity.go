package authentication

import "errors"

type ()

var (
	ErrUserNotFound        = errors.New("invalid username or password")
	ErrInvalidPassword     = errors.New("invalid username or password")
	ErrAuthorizationFailed = errors.New("you are not authorized to access this resource")
)
