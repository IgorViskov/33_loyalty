package core

type Result[T any] struct {
	data *T
	err  error
}

func Done[T any](data *T) Result[T] {
	return Result[T]{data, nil}
}

func Failed[T any](err error) Result[T] {
	return Result[T]{
		err: err,
	}
}

func (r *Result[T]) Success() bool {
	return r.err == nil
}

func (r *Result[T]) Data() *T {
	return r.data
}

func (r *Result[T]) Err() error {
	return r.err
}
