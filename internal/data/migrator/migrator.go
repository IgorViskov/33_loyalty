package migrator

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/data"
	"github.com/IgorViskov/33_loyalty/internal/domain"
)

func AutoMigrate(connector data.Connector) error {
	session := connector.GetConnection(context.Background())
	return session.AutoMigrate(&domain.User{}, &domain.Accrual{}, &domain.AccrualTask{}, domain.Withdrawals{}, domain.PoisonQueueItem{})
}
