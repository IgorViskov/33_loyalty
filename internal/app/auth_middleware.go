package app

import (
	"github.com/IgorViskov/33_loyalty/internal/api"
	"github.com/IgorViskov/33_loyalty/internal/api/models"
	"github.com/IgorViskov/33_loyalty/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		uc := api.GetContext(ec)

		path := uc.Request().URL.Path
		if path == "/api/user/register" || path == "/api/user/login" {
			return next(ec)
		}

		cookie, err := ec.Request().Cookie(config.AuthCookieName)
		if err != nil {
			return ec.NoContent(http.StatusUnauthorized)
		}

		claims := &models.Claims{}
		_, err = jwt.ParseWithClaims(cookie.Value, claims, func(_ *jwt.Token) (interface{}, error) {
			return []byte(config.AuthSecretKey), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		uc.UserService = uc.UserService.Login(claims.UserID)

		return next(ec)
	}
}
