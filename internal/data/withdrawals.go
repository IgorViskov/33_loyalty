package data

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"gorm.io/gorm"
)

type WithdrawalsRepository struct {
	connector Connector
}

func NewWithdrawalsRepository(connector Connector) domain.WithdrawalsRepository {
	return &WithdrawalsRepository{
		connector: connector,
	}
}

func (s *WithdrawalsRepository) Get(context context.Context, id uint64) (*domain.Withdrawals, error) {
	session := s.getSession(context)
	var r domain.Withdrawals
	err := session.First(&r, id).Error
	return &r, err
}

func (s *WithdrawalsRepository) Insert(context context.Context, entity *domain.Withdrawals) (*domain.Withdrawals, error) {
	session := s.getSession(context)
	result := session.Create(entity)
	err := result.Error
	return entity, err
}

func (s *WithdrawalsRepository) Update(_ context.Context, _ *domain.Withdrawals) (*domain.Withdrawals, error) {
	return nil, apperrors.ErrNonImplemented
}

func (s *WithdrawalsRepository) Delete(_ context.Context, _ uint64) error {
	return apperrors.ErrNonImplemented
}

func (s *WithdrawalsRepository) Close() error {
	return s.connector.Close()
}

func (s *WithdrawalsRepository) AllByUser(ctx context.Context, userID uint64) ([]domain.Withdrawals, error) {
	session := s.getSession(ctx)
	entities := make([]domain.Withdrawals, 0)
	result := session.Where(&domain.Withdrawals{UserID: userID}).Find(&entities)
	return entities, result.Error
}

func (s *WithdrawalsRepository) getSession(c context.Context) *gorm.DB {
	return s.connector.GetConnection(c)
}
