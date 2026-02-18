package vulcan

import (
	"context"
)

// Модель можно загрузить по схеме структуры в виде плоского MAP.
func (q *Query[T]) LoadMap(ctx context.Context) ([]map[string]any, map[string][]any, error) {
	mapData := []map[string]any{}
	pkMap := map[string][]any{}
	q.Build()
	rows, err := q.db.QueryContext(ctx, q.fullStatement, q.Bindings...)
	if err != nil {
		return mapData, pkMap, err
	}

	defer rows.Close()

	cols, _ := rows.Columns()

	table := q.Model.TableName
	pk := q.Model.Pk

	pkCols := q.getPk(pk, table)

	n := len(cols)
	colValues := make([]any, n)
	colPtrs := make([]any, n)
	for i := range colValues {
		colPtrs[i] = &colValues[i]
	}

	for rows.Next() {
		if err := rows.Scan(colPtrs...); err != nil {
			panic(err)
		}

		colsMap := make(map[string]any, n)
		for i, col := range cols {
			colsMap[col] = colValues[i]
		}

		for _, pkCol := range pkCols {
			_, e := colsMap[pkCol]
			if !e {
				continue
			}
			pkMap[pkCol] = append(pkMap[pkCol], colsMap[pkCol])
		}

		mapData = append(mapData, colsMap)
	}

	return mapData, pkMap, nil
}

func (q *Query[T]) CLoad(ctx context.Context) ([]T, error) {
	var model T
	parentData, parentPkMap, err := q.LoadMap(ctx)
	if err != nil {
		return []T{}, err
	}
	data, err := q.smartHydration(ctx, &model, parentData, parentPkMap)
	if err != nil {
		return nil, err
	}
	return q.reflectSliceToSlice(data), nil
}

func (q *Query[T]) Load(ctx context.Context) ([]T, error) {
	var model T
	parentData, parentPkMap, err := q.LoadMap(ctx)
	if err != nil {
		return []T{}, err
	}
	data, err := q.smartHydrationSync(ctx, &model, parentData, parentPkMap)
	if err != nil {
		return nil, err
	}
	return q.reflectSliceToSlice(data), nil
}

func (q *Query[T]) With(field string, closure func(*Query[T])) *Query[T] {
	if q.whereHasMap == nil {
		q.whereHasMap = make(map[string]func(*Query[T]))
	}
	q.whereHasMap[field] = closure
	return q
}
