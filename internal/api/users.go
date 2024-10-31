package api

import (
	"errors"
	"github.com/IgorViskov/33_loyalty/internal/api/models"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Login(c echo.Context) error {
	uc := GetContext(c)
	var model models.AuthModel
	err := c.Bind(&model)
	if err != nil {
		return c.String(http.StatusBadRequest, apperrors.ErrInvalidFormatRequest.Error())
	}
	return authenticate(uc, model)
}

func Register(c echo.Context) error {
	uc := GetContext(c)
	var model models.AuthModel
	err := c.Bind(&model)
	if err != nil {
		return c.String(http.StatusBadRequest, apperrors.ErrInvalidFormatRequest.Error())
	}
	err = uc.UserService.Register(uc.GetCtx(), model.Login, model.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrInsertConflict) {
			return uc.NoContent(http.StatusConflict)
		} else {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return authenticate(uc, model)
}

func authenticate(uc *UserContext, model models.AuthModel) error {
	user, err := uc.UserService.CheckPassword(uc.GetCtx(), model.Login, model.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrRecordNotFound) {
			return uc.String(http.StatusUnauthorized, apperrors.ErrPairLoginPasswordNotValid.Error())
		} else {
			return uc.String(http.StatusUnauthorized, err.Error())
		}

	}
	err = saveCookie(uc.Response().Writer, user.ID)
	if err != nil {
		return uc.String(http.StatusInternalServerError, err.Error())
	}
	return uc.NoContent(http.StatusOK)
}

func saveCookie(w http.ResponseWriter, userID uint64) error {
	cookie, err := createCookie(&models.Claims{
		UserID: userID,
	})
	if err != nil {
		return err
	}

	http.SetCookie(w, cookie)
	return nil
}

func createCookie(claims *models.Claims) (*http.Cookie, error) {
	val, err := getToken(claims)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:  config.AuthCookieName,
		Value: val,
		Path:  "/",
	}, nil
}

func getToken(claims *models.Claims) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.AuthSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
