package httpclient

import "errors"

var (
	ErrRequestTimeout  = errors.New("request timeout")
	ErrConnectionError = errors.New("connection error")
)
