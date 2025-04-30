package authenticationUsecase

import (
	"context"

	"github.com/extizout/pokemon-api-server/config"
	"github.com/extizout/pokemon-api-server/pkg/authentication"
	"github.com/extizout/pokemon-api-server/pkg/authentication/authenticationRepository"
	"github.com/extizout/pokemon-api-server/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthenticationUsecaseService interface {
		Login(pctx context.Context, username, password string, secret string) (*string, error)
		VerifyUser(pctx context.Context, userame string) (bool, error)
	}
	authenticationUsecase struct {
		cfg                      *config.Config
		authenticationRepository authenticationRepository.AuthenticationRepositoryService
	}
)

func NewAuthenticationUsecase(cfg *config.Config, authenticationRepository authenticationRepository.AuthenticationRepositoryService) AuthenticationUsecaseService {
	return &authenticationUsecase{
		cfg:                      cfg,
		authenticationRepository: authenticationRepository,
	}
}

func (a *authenticationUsecase) Login(pctx context.Context, username, password string, secret string) (*string, error) {
	result, err := a.authenticationRepository.FindOneUserAndCheckCredential(pctx, username, password)
	if err != nil {
		return nil, authentication.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.HashedPassword), []byte(password)); err != nil {
		return nil, authentication.ErrInvalidPassword
	}

	tokenFactory := jwt.NewAccessToken(a.cfg.Jwt.AccessSecretKey, a.cfg.Jwt.AccessDuration, &authentication.Claims{
		Username: result.Username,
	})
	token := tokenFactory.SignToken()

	return &token, nil

}

func (a *authenticationUsecase) VerifyUser(pctx context.Context, userame string) (bool, error) {
	_, err := a.authenticationRepository.FindOneUserAndCheckCredential(pctx, userame, "")
	if err != nil {
		return false, authentication.ErrAuthorizationFailed
	}

	return true, nil
}
