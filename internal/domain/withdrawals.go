package domain

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/core"
	"github.com/govalues/decimal"
	"time"
)

type Withdrawals struct {
	ID          uint64              `gorm:"primary_key"`
	Value       decimal.NullDecimal `gorm:"type: numeric"`
	SpentAt     time.Time
	UserID      uint64 `gorm:"index"`
	OrderNumber string `gorm:"index;unique;size:15"`
}

type WithdrawalsRepository interface {
	core.Repository[uint64, Withdrawals]
	AllByUser(ctx context.Context, userID uint64) ([]Withdrawals, error)
}
