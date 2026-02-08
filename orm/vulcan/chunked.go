package vulcan

import "context"

// IN WORK!
func (q *Query[T]) Chunk(ctx context.Context, chunk int, closure func([]T) error) error {
	var prev any
	prev = nil
	for {
		data, err := q.CursorPaginate("id", prev, chunk).Load(ctx)
		prev = data[len(data)-1] // TODO: REFLECT GET ID FROM METADATA!

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
