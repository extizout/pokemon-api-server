package server

import (
	"github.com/extizout/pokemon-api-server/modules/user"
	"github.com/extizout/pokemon-api-server/modules/user/userHandler"
	"github.com/extizout/pokemon-api-server/modules/user/userRepository"
	"github.com/extizout/pokemon-api-server/modules/user/userUsecase"
)

func (s *server) userService(usersInMemory *[]user.User) {
	repository := userRepository.NewUserRepository(usersInMemory)
	usecase := userUsecase.NewUserUsecase(repository)
	handler := userHandler.NewUserHttpHandler(s.cfg, usecase)

	userRouter := s.app.Group("/api/v1/user")

	userRouter.POST("/register", handler.CreateUser)
}
