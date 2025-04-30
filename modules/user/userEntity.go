package user

import (
	"errors"
	"time"
)

type (
	User struct {
		Id             int       `json:"id"`
		Username       string    `json:"username"`
		HashedPassword string    `json:"hashedPassword"`
		CreatedAt      time.Time `json:"createdAt"`
		UpdatedAt      time.Time `json:"updatedAt"`
	}
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrHashPasswordFailed = errors.New("failed to hash password")
	ErrCountUsersFailed   = errors.New("failed to count users")
	ErrCreateUserFailed   = errors.New("failed to create user")
)
