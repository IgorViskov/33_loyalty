package api

import (
	"github.com/IgorViskov/33_loyalty/internal/api/models"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/services"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type Controller struct {
	service   *services.OrdersService
	withdraws *services.WithdrawService
}

func NewController(service *services.OrdersService, withdraws *services.WithdrawService) *Controller {
	return &Controller{
		service:   service,
		withdraws: withdraws,
	}
}

func (c *Controller) RegisterOrder(ec echo.Context) error {
	uc := GetContext(ec)
	body, err := io.ReadAll(ec.Request().Body)
	if err != nil {
		return uc.String(http.StatusBadRequest, err.Error())
	}
	orderNumber := string(body)
	if err = goluhn.Validate(orderNumber); err != nil {
		return uc.String(http.StatusUnprocessableEntity, apperrors.MsgIncorrectOrderNumber)
	}

	a, err := c.service.GetByOrder(uc.GetCtx(), orderNumber)
	if err != nil {
		return uc.String(http.StatusInternalServerError, err.Error())
	}
	if a != nil {
		if a.UserID == *uc.UserService.GetUserID() {
			return uc.NoContent(http.StatusOK)
		}

		return uc.String(http.StatusConflict, apperrors.MsgOrderEasUploadedAnotherUser)
	}

	_ = c.service.Enqueue(uc.GetCtx(), orderNumber, *uc.UserService.GetUserID())

	return uc.NoContent(http.StatusAccepted)
}

func (c *Controller) GetAllRegisteredOrders(ec echo.Context) error {
	uc := GetContext(ec)
	accruals, err := c.service.GetAll(uc.GetCtx(), *uc.UserService.GetUserID())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if len(accruals) == 0 {
		return uc.NoContent(http.StatusNoContent)
	}

	return uc.JSON(http.StatusOK, accruals)
}

func (c *Controller) Balance(ec echo.Context) error {
	uc := GetContext(ec)
	result, err := c.service.GetBalance(uc.GetCtx(), *uc.UserService.GetUserID())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return uc.JSON(http.StatusOK, &result)
}

func (c *Controller) Withdraw(ec echo.Context) error {
	uc := GetContext(ec)
	var request models.WithdrawRequest
	err := uc.Bind(&request)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err = goluhn.Validate(request.Order); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, apperrors.MsgIncorrectOrderNumber)
	}

	userID := *uc.UserService.GetUserID()
	ctx := uc.GetCtx()

	exist, err := c.service.Exist(ctx, request.Order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if exist {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, apperrors.MsgIncorrectOrderNumber)
	}

	balance, err := c.service.GetBalance(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if balance.Current.Less(request.Sum) {
		return echo.NewHTTPError(http.StatusPaymentRequired, "there are insufficient funds in the account")
	}

	err = c.withdraws.Withdraw(ctx, request.Order, request.Sum, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return uc.NoContent(http.StatusOK)
}

func (c *Controller) AllWithdraw(ec echo.Context) error {
	uc := GetContext(ec)

	response, err := c.withdraws.GetAll(uc.GetCtx(), *uc.UserService.GetUserID())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return uc.JSON(http.StatusOK, response)
}
