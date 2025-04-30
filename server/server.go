package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/extizout/pokemon-api-server/config"
	"github.com/extizout/pokemon-api-server/modules/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type (
	server struct {
		app *echo.Echo
		cfg *config.Config
	}
)

func (s *server) gracefulShutdown(pctx context.Context, quit <-chan os.Signal) {
	<-quit
	log.Warnf("Shutting down server")
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()
	if err := s.app.Shutdown(ctx); err != nil {
		log.Errorf("error: %v", err.Error())
	}
}

func (s *server) httpListening() {
	if err := s.app.Start(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
		log.Errorf("error: %v", err.Error())
	}
}

func Start(pctx context.Context, cfg *config.Config) {
	s := &server{
		app: echo.New(),
		cfg: cfg,
	}

	switch cfg.App.Stage {
	case "development":
		s.app.Debug = true
		s.app.Logger.SetLevel(log.DEBUG)
	default:
		s.app.Logger.SetLevel(log.INFO)
	}

	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	s.app.Use(middleware.BodyLimit(cfg.App.BodyLimit))

	loggerConfig := middleware.LoggerConfig{}

	s.app.Use(middleware.LoggerWithConfig(loggerConfig))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go s.gracefulShutdown(pctx, quit)

	s.app.GET("/health", s.healthCheckService)

	//TODO: decrease coupling btw auth & user service
	users := []user.User{}

	s.authenticateService(&users)
	s.userService(&users)
	s.pokemonService(&users)

	s.httpListening()
}
