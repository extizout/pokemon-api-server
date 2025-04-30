package userHandler

import (
	"context"
	"net/http"

	"github.com/extizout/pokemon-api-server/config"
	"github.com/extizout/pokemon-api-server/modules/user"
	"github.com/extizout/pokemon-api-server/modules/user/userUsecase"
	"github.com/extizout/pokemon-api-server/pkg/utils/request"
	"github.com/extizout/pokemon-api-server/pkg/utils/response"
	"github.com/labstack/echo/v4"
)

type (
	UserHttpHandlerService interface {
		CreateUser(c echo.Context) error
	}

	userHttpHandler struct {
		cfg         *config.Config
		userUsecase userUsecase.UserUsecaseService
	}
)

func NewUserHttpHandler(cfg *config.Config, userUsecase userUsecase.UserUsecaseService) UserHttpHandlerService {
	return &userHttpHandler{cfg: cfg, userUsecase: userUsecase}
}

func (h *userHttpHandler) CreateUser(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)
	requestDto := new(user.CreateUserRequestDto)

	if err := wrapper.Bind(requestDto); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	userId, err := h.userUsecase.CreateUser(ctx, requestDto)
	if err != nil {
		return response.ErrorResponse(c, http.StatusConflict, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, &user.CreateUserResponseDto{
		Id: *userId,
	})
}
