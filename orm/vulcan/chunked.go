package vulcan

import (
	"context"
	"reflect"
)

// Experimental feature!
func (q *Query[T]) ChunkById(ctx context.Context, chunk int, closure func([]T) error) error {
	var prev any
	prev = nil
	for {
		data, err := q.Clone().CursorPaginate(q.Model.Pks[0], prev, chunk).Load(ctx)
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

		prevT := data[len(data)-1]
		prev = reflect.ValueOf(prevT).FieldByName(q.Model.PksInStruct[0]).Interface()
	}
	return nil
}

func (q *Query[T]) Each(ctx context.Context, closure func(*T) error) error {
	err := q.ChunkById(ctx, 1000, func(t []T) error {
		for _, elem := range t {
			err := closure(&elem)
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
