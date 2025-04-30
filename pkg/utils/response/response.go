package response

import "github.com/labstack/echo/v4"

type (
	ResponesDto struct {
		Message string `json:"message"`
		Data    any    `json:"data, omitempty"`
	}
)

func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, ResponesDto{Message: message})
}
func SuccessResponse(c echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, ResponesDto{
		Message: "success",
		Data:    data,
	})
}
