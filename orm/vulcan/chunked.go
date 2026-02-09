package vulcan

import (
	"context"
	"reflect"
)

// IN WORK! Acceptable If Id field is defined. Experimental feature!
func (q *Query[T]) ChunkById(ctx context.Context, chunk int, closure func([]T) error) error {
	var prev any
	prev = nil
	for {
		data, err := q.CursorPaginate("id", prev, chunk).Load(ctx)
		prevT := data[len(data)-1] // TODO: REFLECT GET ID FROM METADATA!
		prev = reflect.ValueOf(prevT).Elem().FieldByName("Id").Interface()

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
	err := q.ChunkById(ctx, 1000, func(t []T) error {
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
