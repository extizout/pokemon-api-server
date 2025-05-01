package userUsecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/extizout/pokemon-api-server/modules/user"
	"github.com/extizout/pokemon-api-server/modules/user/userUsecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

// FindOneByUsername implements userRepository.UserRepositoryService.
func (m *MockUserRepository) FindOneByUsername(pctx context.Context, username string) (*user.User, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) IsUsernameUnique(ctx context.Context, username string) bool {
	args := m.Called(ctx, username)
	return args.Bool(0)
}

func (m *MockUserRepository) CountUsers(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := userUsecase.NewUserUsecase(mockRepo)

	req := &user.CreateUserRequestDto{
		Username: "testuser",
		Password: "password123",
	}

	ctx := context.TODO()

	mockRepo.On("IsUsernameUnique", ctx, "testuser").Return(false)
	mockRepo.On("CountUsers", ctx).Return(1, nil)
	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	id, err := usecase.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, 2, *id)

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_UsernameExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := userUsecase.NewUserUsecase(mockRepo)

	req := &user.CreateUserRequestDto{
		Username: "existinguser",
		Password: "password123",
	}

	ctx := context.TODO()

	mockRepo.On("IsUsernameUnique", ctx, "existinguser").Return(true)

	id, err := usecase.CreateUser(ctx, req)

	assert.Nil(t, id)
	assert.ErrorIs(t, err, user.ErrUserAlreadyExists)
}

func TestCreateUser_CountUsersFail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := userUsecase.NewUserUsecase(mockRepo)

	req := &user.CreateUserRequestDto{
		Username: "newuser",
		Password: "password123",
	}

	ctx := context.TODO()

	mockRepo.On("IsUsernameUnique", ctx, "newuser").Return(false)
	mockRepo.On("CountUsers", ctx).Return(0, errors.New("count error"))

	id, err := usecase.CreateUser(ctx, req)

	assert.Nil(t, id)
	assert.ErrorIs(t, err, user.ErrCountUsersFailed)
}

func TestCreateUser_CreateUserFail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := userUsecase.NewUserUsecase(mockRepo)

	req := &user.CreateUserRequestDto{
		Username: "failuser",
		Password: "password123",
	}

	ctx := context.TODO()

	mockRepo.On("IsUsernameUnique", ctx, "failuser").Return(false)
	mockRepo.On("CountUsers", ctx).Return(5, nil)
	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*user.User")).Return(errors.New("create error"))

	id, err := usecase.CreateUser(ctx, req)

	assert.Nil(t, id)
	assert.ErrorIs(t, err, user.ErrCreateUserFailed)
}
