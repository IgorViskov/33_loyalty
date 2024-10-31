package api

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/services"
	"github.com/labstack/echo/v4"
)

type UserContext struct {
	echo.Context
	UserService *services.UserService
}

func (uc *UserContext) GetCtx() context.Context {
	return uc.Request().Context()
}

func GetContext(ec echo.Context) *UserContext {
	return ec.(*UserContext)
}
