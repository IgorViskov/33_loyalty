package domain

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/core"
)

type User struct {
	ID           uint64
	Login        string `gorm:"unique"`
	PasswordHash string
}

type UserRepository interface {
	core.Repository[uint64, User]
	GetByLogin(context context.Context, login string) (*User, error)
}
