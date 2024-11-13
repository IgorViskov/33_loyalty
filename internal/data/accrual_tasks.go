package data

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"gorm.io/gorm"
)

type AccrualTasksRepository struct {
	connector Connector
}

func NewAccrualTasksRepository(connector Connector) *AccrualTasksRepository {
	return &AccrualTasksRepository{
		connector: connector,
	}
}

func (s *AccrualTasksRepository) Get(context context.Context, id uint64) (*domain.AccrualTask, error) {
	session := s.getSession(context)
	var r domain.AccrualTask
	err := session.First(&r, id).Error
	return &r, err
}

func (s *AccrualTasksRepository) Insert(context context.Context, entity *domain.AccrualTask) (*domain.AccrualTask, error) {
	session := s.getSession(context)
	result := session.Create(entity)
	err := result.Error
	return entity, err
}

func (s *AccrualTasksRepository) Update(_ context.Context, _ *domain.AccrualTask) (*domain.AccrualTask, error) {
	return nil, apperrors.ErrNonImplemented
}

func (s *AccrualTasksRepository) Delete(context context.Context, id uint64) error {
	session := s.getSession(context)
	return session.Delete(&domain.AccrualTask{}, id).Error
}

func (s *AccrualTasksRepository) Close() error {
	return s.connector.Close()
}

func (s *AccrualTasksRepository) All(context context.Context) ([]domain.AccrualTask, error) {
	session := s.getSession(context)
	entities := make([]domain.AccrualTask, 0)
	result := session.Find(&entities)
	return entities, result.Error
}

func (s *AccrualTasksRepository) DeleteFromOrder(context context.Context, order string) error {
	session := s.getSession(context)
	result := session.Where("order_number = $1", order).Delete(&domain.AccrualTask{})
	return result.Error
}

func (s *AccrualTasksRepository) getSession(c context.Context) *gorm.DB {
	return s.connector.GetConnection(c)
}
