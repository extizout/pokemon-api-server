package middleware

import (
	"net/http"
	"strings"

	"github.com/extizout/pokemon-api-server/pkg/authentication"
	"github.com/extizout/pokemon-api-server/pkg/authentication/authenticationUsecase"
	"github.com/extizout/pokemon-api-server/pkg/jwt"
	"github.com/extizout/pokemon-api-server/pkg/utils/response"
	"github.com/labstack/echo/v4"
)

func NewJWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			token, err := jwt.ParseToken(secret, tokenStr)
			if err != nil {
				return response.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			}
			c.Set("username", token.Username)

			return next(c)
		}
	}
}
func NewUserVerificationMiddleware(useCase authenticationUsecase.AuthenticationUsecaseService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username, ok := c.Get("username").(string)
			if !ok || username == "" {
				return response.ErrorResponse(c, http.StatusForbidden, authentication.ErrAuthorizationFailed.Error())
			}

			isAccess, err := useCase.VerifyUser(c.Request().Context(), username)
			if !isAccess {
				return response.ErrorResponse(c, http.StatusForbidden, err.Error())
			}

			return next(c)
		}
	}
}
