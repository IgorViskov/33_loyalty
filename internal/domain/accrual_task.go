package domain

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/core"
	"time"
)

type AccrualTask struct {
	OrderNumber string `gorm:"index;unique;size:15"`
	UploadedAt  time.Time
	UserID      uint64 `gorm:"index"`
}

type AccrualTaskRepository interface {
	core.Repository[uint64, AccrualTask]
	DeleteFromOrder(context context.Context, order string) error
	All(context context.Context) ([]AccrualTask, error)
}
