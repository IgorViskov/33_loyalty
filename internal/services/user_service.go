package services

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/apperrors"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	users  domain.UserRepository
	userID *uint64
}

func NewUserService(users domain.UserRepository, userID *uint64) *UserService {
	return &UserService{
		users:  users,
		userID: userID,
	}
}

func (us *UserService) Login(userID uint64) *UserService {
	return &UserService{
		users:  us.users,
		userID: &userID,
	}
}

func (us *UserService) IsAuth() bool {
	return us.userID != nil
}

func (us *UserService) Register(ctx context.Context, login string, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	_, err = us.users.Insert(ctx, &domain.User{
		Login:        login,
		PasswordHash: hash,
	})

	return err
}

func (us *UserService) GetUserID() *uint64 {
	return us.userID
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (us *UserService) CheckPassword(ctx context.Context, login string, password string) (*domain.User, error) {
	user, err := us.users.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrPairLoginPasswordNotValid
	}
	if !checkPassword(password, user.PasswordHash) {
		return nil, apperrors.ErrPairLoginPasswordNotValid
	}
	return user, nil
}

func checkPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
