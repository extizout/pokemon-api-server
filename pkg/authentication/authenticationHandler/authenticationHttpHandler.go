package authenticationHandler

import (
	"context"
	"net/http"

	"github.com/extizout/pokemon-api-server/config"
	"github.com/extizout/pokemon-api-server/pkg/authentication"
	"github.com/extizout/pokemon-api-server/pkg/authentication/authenticationUsecase"
	"github.com/extizout/pokemon-api-server/pkg/utils/request"
	"github.com/extizout/pokemon-api-server/pkg/utils/response"
	"github.com/labstack/echo/v4"
)

type (
	AuthenticationHttpHandlerService interface {
		Login(c echo.Context) error
	}

	authenticationHttpHandler struct {
		cfg                   *config.Config
		authenticationUsecase authenticationUsecase.AuthenticationUsecaseService
	}
)

func NewAuthenticationHttpHandler(cfg *config.Config, authenticationUsecase authenticationUsecase.AuthenticationUsecaseService) AuthenticationHttpHandlerService {
	return &authenticationHttpHandler{cfg: cfg, authenticationUsecase: authenticationUsecase}
}

func (h *authenticationHttpHandler) Login(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)
	requestDto := new(authentication.UserLoginRequestDto)

	if err := wrapper.Bind(requestDto); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	data, err := h.authenticationUsecase.Login(ctx, requestDto.Username, requestDto.Password, h.cfg.Jwt.AccessSecretKey)
	if err != nil {
		return response.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, &authentication.UserLoginResponseDto{
		AccessToken: *data,
	})
}
