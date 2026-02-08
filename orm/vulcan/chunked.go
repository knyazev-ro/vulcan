package vulcan

import "context"

func (q *Query[T]) Chunk(ctx context.Context, chunk int, closure func([]T) error) error {
	for {
		data, err := q.CursorPaginate("id", nil, chunk).Load(ctx)

		if err != nil {
			return err
		}

		if len(data) <= 0 {
			break
		}

		err = closure(data)

		if err != nil {
			return err
		}
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

func (q *Query[T]) Map(ctx context.Context, closure func(T) (T, error)) ([]T, error) {
	result := []T{}
	err := q.Chunk(ctx, 1000, func(t []T) error {
		for _, elem := range t {
			model, err := closure(elem)
			if err != nil {
				return err
			}
			result = append(result, model)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
