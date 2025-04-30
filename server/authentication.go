package server

import (
	"github.com/extizout/pokemon-api-server/modules/user"
	"github.com/extizout/pokemon-api-server/pkg/authentication/authenticationHandler"
	"github.com/extizout/pokemon-api-server/pkg/authentication/authenticationRepository"
	"github.com/extizout/pokemon-api-server/pkg/authentication/authenticationUsecase"
)

func (s *server) authenticateService(usersInMemory *[]user.User) {
	repository := authenticationRepository.NewAuthenticationRepository(usersInMemory)
	usecase := authenticationUsecase.NewAuthenticationUsecase(s.cfg, repository)
	handler := authenticationHandler.NewAuthenticationHttpHandler(s.cfg, usecase)

	authRouter := s.app.Group("/api/v1/auth")
	authRouter.POST("/login", handler.Login)
}
