package userUsecase

import (
	"context"
	"time"

	"github.com/extizout/pokemon-api-server/modules/user"
	"github.com/extizout/pokemon-api-server/modules/user/userRepository"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserUsecaseService interface {
		CreateUser(pctx context.Context, createUserRequestDto *user.CreateUserRequestDto) (*int, error)
	}
	userUsecase struct {
		userRepository userRepository.UserRepositoryService
	}
)

func NewUserUsecase(userRepository userRepository.UserRepositoryService) UserUsecaseService {
	return &userUsecase{userRepository: userRepository}
}

func (u *userUsecase) CreateUser(pctx context.Context, createUserRequestDto *user.CreateUserRequestDto) (*int, error) {
	if u.userRepository.IsUsernameUnique(pctx, createUserRequestDto.Username) {
		return nil, user.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserRequestDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, user.ErrHashPasswordFailed
	}

	count, err := u.userRepository.CountUsers(pctx)
	if err != nil {
		return nil, user.ErrCountUsersFailed
	}

	newUser := &user.User{
		Id:             count + 1,
		Username:       createUserRequestDto.Username,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = u.userRepository.CreateUser(pctx, newUser)
	if err != nil {
		return nil, user.ErrCreateUserFailed
	}

	return &newUser.Id, nil
}
