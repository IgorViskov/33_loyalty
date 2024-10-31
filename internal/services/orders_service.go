package services

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/api/models"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"github.com/IgorViskov/33_loyalty/internal/domain/statuses"
	"github.com/govalues/decimal"
	"time"
)

type OrdersService struct {
	accruals    domain.AccrualRepository
	withdrawals domain.WithdrawalsRepository
	taskService *AccrualTasksService
}

func NewOrdersService(accruals domain.AccrualRepository, tasks *AccrualTasksService, withdrawals domain.WithdrawalsRepository) *OrdersService {
	return &OrdersService{
		accruals:    accruals,
		taskService: tasks,
		withdrawals: withdrawals,
	}
}

func (o *OrdersService) GetByOrder(context context.Context, order string) (*domain.Accrual, error) {
	return o.accruals.FindOrder(context, order)
}

func (o *OrdersService) Enqueue(context context.Context, order string, userID uint64) error {
	err := o.taskService.Enqueue(domain.AccrualTask{
		OrderNumber: order,
		UserID:      userID,
		UploadedAt:  time.Now(),
	})
	if err != nil {
		return err
	}

	_, err = o.accruals.Insert(context, &domain.Accrual{
		OrderNumber: order,
		UserID:      userID,
		UploadedAt:  time.Now(),
		Status:      statuses.NEW,
	})
	return err
}

func (o *OrdersService) GetAll(context context.Context, userID uint64) ([]models.OrdersResponse, error) {
	acc, err := o.accruals.AllByUser(context, userID)
	if err != nil {
		return nil, err
	}
	return mapper(acc), nil
}

func mapper(accruals []domain.Accrual) []models.OrdersResponse {
	res := make([]models.OrdersResponse, len(accruals))
	for i, acc := range accruals {
		res[i] = models.OrdersResponse{
			Accrual:    acc.Value,
			Status:     acc.Status,
			UploadedAt: acc.UploadedAt,
			Number:     acc.OrderNumber,
		}
	}
	return res
}

func (o *OrdersService) GetBalance(context context.Context, userID uint64) (models.BalanceResponse, error) {
	accruals, err := o.accruals.AllByUser(context, userID)
	if err != nil {
		return models.BalanceResponse{}, err
	}
	withdrawals, err := o.withdrawals.AllByUser(context, userID)
	if err != nil {
		return models.BalanceResponse{}, err
	}
	result := decimal.Zero
	for _, a := range accruals {
		if !a.Value.Valid {
			continue
		}
		result, err = result.Add(a.Value.Decimal)
		if err != nil {
			return models.BalanceResponse{}, err
		}
	}
	withdrawalsResult := decimal.Zero
	for _, w := range withdrawals {
		if !w.Value.Valid {
			continue
		}
		withdrawalsResult, err = withdrawalsResult.Add(w.Value.Decimal)
		if err != nil {
			return models.BalanceResponse{}, err
		}
	}

	result, err = result.Sub(withdrawalsResult)

	return models.BalanceResponse{
		Current:   result,
		Withdrawn: withdrawalsResult,
	}, err
}

func (o *OrdersService) Exist(context context.Context, order string) (bool, error) {
	a, e := o.accruals.FindOrder(context, order)
	if e != nil {
		return false, e
	}
	return a != nil, nil
}
