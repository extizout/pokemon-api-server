package server

import (
	"net/http"

	"github.com/extizout/pokemon-api-server/pkg/utils/response"
	"github.com/labstack/echo/v4"
)

type (
	healthCheckDto struct {
		AppName string `json:"appName"`
		Version string `json:"version"`
		Status  string `json:"status"`
	}
)

func (s *server) healthCheckService(c echo.Context) error {
	health := &healthCheckDto{
		AppName: "Pokemon API",
		Version: "1.0.0",
		Status:  "healthy",
	}
	return response.SuccessResponse(c, http.StatusOK, health)
}
