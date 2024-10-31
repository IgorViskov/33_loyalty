package data

import (
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type PoisonRepository struct {
	connector Connector
}

func NewPoisonRepository(connector Connector) *AccrualRepository {
	return &AccrualRepository{
		connector: connector,
	}
}

func (s *PoisonRepository) Get(context context.Context, id uint64) (*domain.PoisonQueueItem, error) {
	session := s.getSession(context)
	var r domain.PoisonQueueItem
	err := session.First(&r, id).Error
	return &r, err
}

func (s *PoisonRepository) Insert(context context.Context, entity *domain.PoisonQueueItem) (*domain.PoisonQueueItem, error) {
	session := s.getSession(context)
	result := session.Create(entity)
	err := result.Error
	return entity, err
}

func (s *PoisonRepository) Update(_ context.Context, _ *domain.PoisonQueueItem) (*domain.PoisonQueueItem, error) {
	return nil, apperrors.ErrNonImplemented
}

func (s *PoisonRepository) Delete(_ context.Context, _ uint64) error {
	return apperrors.ErrNonImplemented
}

func (s *PoisonRepository) Close() error {
	return s.connector.Close()
}

func (s *PoisonRepository) getSession(c context.Context) *gorm.DB {
	return s.connector.GetConnection(c)
}
