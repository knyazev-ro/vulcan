package vulcan

import "context"

func (q *Query[T]) Chunk(ctx context.Context, chunk int, closure func([]T) error) error {
	arr := []T{}
	err := closure(arr)
	if err != nil {
		return err
	}
	return nil
}

func (q *Query[T]) Each(ctx context.Context, closure func(T) error) error {
	err := q.Chunk(ctx, 1000, func(t []T) error {
		for _, elem := range t {
			err := closure(elem)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
