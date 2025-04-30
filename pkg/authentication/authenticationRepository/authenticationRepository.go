package authenticationRepository

import (
	"context"
	"time"

	"github.com/extizout/pokemon-api-server/modules/user"
	"github.com/extizout/pokemon-api-server/pkg/authentication"
)

type AuthenticationRepositoryService interface {
	FindOneUserAndCheckCredential(pctx context.Context, username, password string) (*user.User, error)
}

type authenticationRepository struct {
	users *[]user.User
}

func NewAuthenticationRepository(users *[]user.User) AuthenticationRepositoryService {
	return &authenticationRepository{
		users: users,
	}
}

func (a *authenticationRepository) FindOneUserAndCheckCredential(pctx context.Context, username, password string) (*user.User, error) {
	_, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	for _, u := range *a.users {
		if u.Username == username {
			return &u, nil
		}
	}

	return nil, authentication.ErrUserNotFound
}
