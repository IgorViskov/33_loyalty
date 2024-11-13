package services

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"github.com/govalues/decimal"
	"time"
)

type AccountService struct {
	repo domain.WithdrawalsRepository
}

func NewAccountService(repo domain.WithdrawalsRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (as *AccountService) Add(value decimal.Decimal, userID uint64) error {
	if value.IsNeg() {
		return apperrors.ErrWithdrawalsNegative
	}
	_, err := as.repo.Insert(context.Background(), &domain.Withdrawals{
		Value: decimal.NullDecimal{
			Decimal: value,
			Valid:   true,
		},
		UserID:  userID,
		SpentAt: time.Now(),
	})

	return err
}
