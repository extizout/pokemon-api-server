package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type (
	contextWrapperService interface {
		Bind(data any) error
	}

	contextWrapper struct {
		context   echo.Context
		validator *validator.Validate
	}
)

func ContextWrapper(ctx echo.Context) contextWrapperService {
	return &contextWrapper{
		context:   ctx,
		validator: validator.New(),
	}
}

func (c *contextWrapper) Bind(data any) error {
	if err := c.context.Bind(data); err != nil {
		log.Warnf("failed to bind request data: %v", err)
		return err
	}
	if err := c.validator.Struct(data); err != nil {
		log.Warnf("failed to validate request data: %v", err)
		return err
	}

	return nil
}
