package query

import (
	"fmt"
	"reflect"
)

func (q *Query[T]) reflectSliceToSlice(v []reflect.Value) []T {
	data := make([]T, len(v))
	for i, val := range v {
		data[i] = val.Interface().(T)
	}
	return data
}

func (q *Query[T]) recHydration(data []map[string]any, i reflect.Value) []reflect.Value {
	modelMetadata, ok := i.Type().Elem().FieldByName("_")

	if !ok {
		panic("model metadata was not found")
	}

	pk := modelMetadata.Tag.Get("pk")
	tableName := modelMetadata.Tag.Get("table")

	mapByPk := map[int64][]map[string]any{}
	for _, row := range data {
		pKeyCol := fmt.Sprintf("%s_%s", tableName, pk)
		pKeyVal := row[pKeyCol].(int64)
		mapByPk[pKeyVal] = append(mapByPk[pKeyVal], row)
	}

	structData := []reflect.Value{}

	for _, rowsByPk := range mapByPk {
		newModel := reflect.New(i.Type().Elem()).Elem()
		valModelType := newModel.Type()
		first := rowsByPk[0]
		for i := range newModel.NumField() {
			field := newModel.Field(i)
			fieldType := valModelType.Field(i)
			tagType := fieldType.Tag.Get("type")

			if tagType == "column" {
				col := fieldType.Tag.Get("col")
				colKey := fmt.Sprintf("%s_%s", tableName, col)

				// int64
				if field.Kind() == reflect.Int64 {
					field.SetInt(first[colKey].(int64))
				}
				// string
				if field.Kind() == reflect.String {
					field.SetString(first[colKey].(string))
				}

			}

			if tagType == "relation" {
				relTypeTag := fieldType.Tag.Get("reltype")
				if relTypeTag == OneToMany && field.Kind() == reflect.Slice {
					sliceType := field.Type()    // []Some
					elemType := sliceType.Elem() // Some
					childrens := q.recHydration(rowsByPk, reflect.New(elemType))
					newSlice := reflect.MakeSlice(sliceType, 0, len(childrens))
					for _, child := range childrens {
						newSlice = reflect.Append(newSlice, child)
					}
					field.Set(newSlice)
				}
			}
		}
		structData = append(structData, newModel)
	}
	return structData
}

func (q *Query[T]) HydrationOneToOne(data []map[string]any) []T {

	return []T{}
}

func (q *Query[T]) HydrationOneToMany(data []map[string]any) []T {
	var m T
	total := q.recHydration(data, reflect.ValueOf(&m))
	return q.reflectSliceToSlice(total)
}

func (q *Query[T]) HydrationManyToMany(data []map[string]any) []T {

	return []T{}
}
