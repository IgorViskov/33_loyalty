package services

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/api/models"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"github.com/govalues/decimal"
	"time"
)

type WithdrawService struct {
	repo domain.WithdrawalsRepository
}

func NewWithdrawService(repo domain.WithdrawalsRepository) *WithdrawService {
	return &WithdrawService{repo: repo}
}

func (w *WithdrawService) Withdraw(ctx context.Context, order string, amount decimal.Decimal, userID uint64) error {
	if amount.IsNeg() || amount.IsZero() {
		return apperrors.ErrAmountNotPositive
	}

	_, err := w.repo.Insert(ctx, &domain.Withdrawals{
		Value: decimal.NullDecimal{
			Decimal: amount,
			Valid:   true,
		},
		UserID:      userID,
		SpentAt:     time.Now(),
		OrderNumber: order,
	})

	return err
}

func (w *WithdrawService) GetAll(ctx context.Context, userID uint64) ([]models.WithdrawResponse, error) {
	wi, err := w.repo.AllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapperWithdrawal(wi), nil
}

func mapperWithdrawal(withdrawals []domain.Withdrawals) []models.WithdrawResponse {
	result := make([]models.WithdrawResponse, len(withdrawals))
	for i, w := range withdrawals {
		result[i] = models.WithdrawResponse{
			Order:       w.OrderNumber,
			Sum:         w.Value.Decimal,
			ProcessedAt: w.SpentAt,
		}
	}
	return result
}
