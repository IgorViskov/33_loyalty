package data

import (
	"context"
	"errors"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccrualRepository struct {
	connector Connector
}

func NewAccrualRepository(connector Connector) *AccrualRepository {
	return &AccrualRepository{
		connector: connector,
	}
}

func (s *AccrualRepository) Get(context context.Context, id uint64) (*domain.Accrual, error) {
	session := s.getSession(context)
	var r domain.Accrual
	err := session.First(&r, id).Error
	return &r, err
}

func (s *AccrualRepository) Insert(context context.Context, entity *domain.Accrual) (*domain.Accrual, error) {
	session := s.getSession(context)
	result := session.Create(entity)
	err := result.Error
	return entity, err
}

func (s *AccrualRepository) Update(_ context.Context, _ *domain.Accrual) (*domain.Accrual, error) {
	return nil, apperrors.ErrNonImplemented
}

func (s *AccrualRepository) Delete(_ context.Context, _ uint64) error {
	return apperrors.ErrNonImplemented
}

func (s *AccrualRepository) Close() error {
	return s.connector.Close()
}

func (s *AccrualRepository) All(context context.Context) ([]domain.Accrual, error) {
	session := s.getSession(context)
	entities := make([]domain.Accrual, 0)
	result := session.Find(&entities)
	return entities, result.Error
}

func (s *AccrualRepository) FindOrder(context context.Context, order string) (*domain.Accrual, error) {
	session := s.getSession(context)
	r := &domain.Accrual{
		OrderNumber: order,
	}
	err := session.Where(r).First(r).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return r, err
}

func (s *AccrualRepository) CreateOrUpdate(context context.Context, a *domain.Accrual) (*domain.Accrual, error) {
	session := s.getSession(context)
	result := session.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "order_number"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "value"}),
	}).Create(a)
	err := result.Error
	return a, err
}

func (s *AccrualRepository) AllByUser(context context.Context, userID uint64) ([]domain.Accrual, error) {
	session := s.getSession(context)
	entities := make([]domain.Accrual, 0)
	result := session.Where(&domain.Accrual{UserID: userID}).Find(&entities)
	return entities, result.Error
}

func (s *AccrualRepository) getSession(c context.Context) *gorm.DB {
	return s.connector.GetConnection(c)
}
