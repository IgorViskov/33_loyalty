package domain

import (
	"github.com/IgorViskov/33_loyalty/internal/core"
	"time"
)

type PoisonQueueItem struct {
	ID          uint64 `gorm:"primary_key"`
	OrderNumber uint64 `gorm:"unique"`
	UploadedAt  time.Time
	LastTry     time.Time
	Error       string
	UserID      uint64 `gorm:"index"`
}

type PoisonQueueRepository interface {
	core.Repository[uint64, PoisonQueueItem]
}
