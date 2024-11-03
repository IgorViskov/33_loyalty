package domain

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/core"
	"github.com/IgorViskov/33_loyalty/internal/domain/statuses"
	"github.com/govalues/decimal"
	"time"
)

type Accrual struct {
	ID          uint64 `gorm:"primary_key"`
	OrderNumber string `gorm:"index;unique;size:15"`
	Status      statuses.ProcessStatus
	Value       decimal.NullDecimal `gorm:"type: numeric"`
	UploadedAt  time.Time
	UserID      uint64 `gorm:"index"`
}

type AccrualRepository interface {
	core.Repository[uint64, Accrual]
	All(context context.Context) ([]Accrual, error)
	FindOrder(context context.Context, order string) (*Accrual, error)
	CreateOrUpdate(context context.Context, accrual *Accrual) (*Accrual, error)
	AllByUser(context context.Context, userID uint64) ([]Accrual, error)
}
