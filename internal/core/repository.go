package core

import "context"

type Repository[TKey comparable, TEntity any] interface {
	Get(context context.Context, id TKey) (*TEntity, error)
	Insert(context context.Context, entity *TEntity) (*TEntity, error)
	Update(context context.Context, entity *TEntity) (*TEntity, error)
	Delete(context context.Context, id TKey) error
	Close() error
}
