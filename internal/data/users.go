package data

import (
	"context"
	"errors"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UsersRepository struct {
	connector Connector
}

func NewUsersRepository(connector Connector) *UsersRepository {
	return &UsersRepository{
		connector: connector,
	}
}

func (s *UsersRepository) Get(context context.Context, id uint64) (*domain.User, error) {
	session := s.getSession(context)
	var r domain.User
	err := session.First(&r, id).Error
	return &r, err
}

func (s *UsersRepository) Insert(context context.Context, entity *domain.User) (*domain.User, error) {
	session := s.getSession(context)
	result := session.Clauses(clause.OnConflict{DoNothing: true}).Create(entity)
	if result.RowsAffected == 0 {
		if result.Error == nil {
			return nil, apperrors.ErrInsertConflict
		} else {
			return nil, result.Error
		}
	}
	return entity, nil
}

func (s *UsersRepository) Update(_ context.Context, _ *domain.User) (*domain.User, error) {
	return nil, apperrors.ErrNonImplemented
}

func (s *UsersRepository) Delete(_ context.Context, _ uint64) error {
	return apperrors.ErrNonImplemented
}

func (s *UsersRepository) Close() error {
	return s.connector.Close()
}

func (s *UsersRepository) GetByLogin(context context.Context, login string) (*domain.User, error) {
	session := s.getSession(context)
	var user domain.User
	err := session.Where("login = $1", login).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrRecordNotFound
	}
	return &user, err
}

func (s *UsersRepository) getSession(c context.Context) *gorm.DB {
	return s.connector.GetConnection(c)
}
