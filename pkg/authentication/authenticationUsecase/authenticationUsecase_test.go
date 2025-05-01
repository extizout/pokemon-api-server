package authenticationUsecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/extizout/pokemon-api-server/config"
	"github.com/extizout/pokemon-api-server/modules/user"
	"github.com/extizout/pokemon-api-server/pkg/authentication"
	"github.com/extizout/pokemon-api-server/pkg/authentication/authenticationUsecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) FindOneUserAndCheckCredential(ctx context.Context, username, password string) (*user.User, error) {
	args := m.Called(ctx, username, password)
	if obj := args.Get(0); obj != nil {
		return obj.(*user.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessSecretKey: "testsecret",
			AccessDuration:  60,
		},
	}

	usecase := authenticationUsecase.NewAuthenticationUsecase(cfg, mockRepo)

	user := &user.User{
		Username:       "testuser",
		HashedPassword: "$2a$10$fdaAciH5IK5xnq6mp5Z90Oe2ZG7t1/r40QHAnXOzZ6E55NIYe1nTW",
	}

	ctx := context.TODO()
	mockRepo.On("FindOneUserAndCheckCredential", ctx, "testuser", "1234").Return(user, nil)

	token, err := usecase.Login(ctx, "testuser", "1234", cfg.Jwt.AccessSecretKey)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessSecretKey: "testsecret",
			AccessDuration:  60 * 15,
		},
	}

	usecase := authenticationUsecase.NewAuthenticationUsecase(cfg, mockRepo)

	ctx := context.TODO()
	mockRepo.On("FindOneUserAndCheckCredential", ctx, "unknown", "password").Return(nil, errors.New("not found"))

	token, err := usecase.Login(ctx, "unknown", "password", cfg.Jwt.AccessSecretKey)

	assert.Nil(t, token)
	assert.ErrorIs(t, err, authentication.ErrUserNotFound)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessSecretKey: "testsecret",
			AccessDuration:  60 * 15,
		},
	}

	usecase := authenticationUsecase.NewAuthenticationUsecase(cfg, mockRepo)

	user := &user.User{
		Username:       "testuser",
		HashedPassword: "$2a$10$incorrecthashshouldfail",
	}

	ctx := context.TODO()
	mockRepo.On("FindOneUserAndCheckCredential", ctx, "testuser", "wrongpass").Return(user, nil)

	token, err := usecase.Login(ctx, "testuser", "wrongpass", cfg.Jwt.AccessSecretKey)

	assert.Nil(t, token)
	assert.ErrorIs(t, err, authentication.ErrInvalidPassword)
}

func TestVerifyUser_Success(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	cfg := &config.Config{}

	usecase := authenticationUsecase.NewAuthenticationUsecase(cfg, mockRepo)

	user := &user.User{
		Username: "verifyme",
	}

	ctx := context.TODO()
	mockRepo.On("FindOneUserAndCheckCredential", ctx, "verifyme", "").Return(user, nil)

	result, err := usecase.VerifyUser(ctx, "verifyme")

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestVerifyUser_Failure(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	cfg := &config.Config{}

	usecase := authenticationUsecase.NewAuthenticationUsecase(cfg, mockRepo)

	ctx := context.TODO()
	mockRepo.On("FindOneUserAndCheckCredential", ctx, "failuser", "").Return(nil, errors.New("not found"))

	result, err := usecase.VerifyUser(ctx, "failuser")

	assert.ErrorIs(t, err, authentication.ErrAuthorizationFailed)
	assert.False(t, result)
}
