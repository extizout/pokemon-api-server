package userRepository

import (
	"context"
	"time"

	"github.com/extizout/pokemon-api-server/modules/user"
)

type UserRepositoryService interface {
	IsUsernameUnique(pctx context.Context, username string) bool
	CreateUser(pctx context.Context, user *user.User) error
	CountUsers(pctx context.Context) (int, error)
	FindOneByUsername(pctx context.Context, username string) (*user.User, error)
}

type userRepository struct {
	db *[]user.User
}

func NewUserRepository(db *[]user.User) UserRepositoryService {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) IsUsernameUnique(pctx context.Context, username string) bool {
	_, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	for _, u := range *r.db {
		if u.Username == username {
			return true
		}
	}

	return false
}

func (r *userRepository) CreateUser(pctx context.Context, user *user.User) error {
	_, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	*r.db = append(*r.db, *user)

	return nil
}

func (r *userRepository) FindOneByUsername(pctx context.Context, username string) (*user.User, error) {
	_, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	for _, u := range *r.db {
		if u.Username == username {
			return &u, nil
		}
	}

	return nil, user.ErrUserNotFound
}

func (r *userRepository) CountUsers(pctx context.Context) (int, error) {
	_, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	count := len(*r.db)

	return count, nil
}
