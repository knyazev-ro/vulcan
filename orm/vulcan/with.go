package vulcan

import (
	"fmt"
	"reflect"

	"github.com/knyazev-ro/vulcan/orm/consts"
	"github.com/knyazev-ro/vulcan/orm/db"
)

func MapByKey(data []map[string]any, key string) map[any][]map[string]any {
	result := make(map[any][]map[string]any)
	for _, item := range data {
		if val, exists := item[key]; exists {
			result[val] = append(result[val], item)
		}
	}
	return result
}

func (q *Query[T]) getMetadata() (string, string) {
	var model T
	val := reflect.ValueOf(&model).Elem()
	if val.Kind() != reflect.Struct {
		panic("getMetadata: T must be a struct type")
	}
	metadata, ok := val.Type().FieldByName("_")

	if !ok {
		panic("getMetadata: struct T must have a metadata field")
	}

	typeName := metadata.Tag.Get("type")
	if typeName != "metadata" {
		panic("getMetadata: struct T must have a metadata field with type:metadata tag")
	}

	table := metadata.Tag.Get("table")
	pk := metadata.Tag.Get("pk")
	return table, pk
}

func (q *Query[T]) LoadMap() ([]map[string]any, map[string][]any) {
	q.Build()
	db := db.DB
	println(q.SQL())
	rows, err := db.Query(q.fullStatement, q.Bindings...)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	cols, _ := rows.Columns()
	table, pk := q.getMetadata()
	pkCols := q.getPk(pk, table)
	mapData := []map[string]any{}
	pkMap := map[string][]any{}
	for rows.Next() {
		colValues := make([]any, len(cols))
		colPtrs := make([]any, len(cols))
		colsMap := map[string]any{}

		for i := range colValues {
			colPtrs[i] = &colValues[i]
		}
		if err := rows.Scan(colPtrs...); err != nil {
			panic(err)
		}

		for i, col := range cols {
			colsMap[col] = colValues[i]
		}

		for _, pkCol := range pkCols {
			pkMap[pkCol] = append(pkMap[pkCol], colsMap[pkCol])
		}

		mapData = append(mapData, colsMap)
	}
	return mapData, pkMap
}

func (q *Query[T]) fillCols(i interface{}, row map[string]any, TableName string) any {
	val := reflect.ValueOf(i).Elem()
	if val.Kind() != reflect.Struct {
		panic("fillCols: i must be a pointer to a struct")
	}

	metadata, ok := val.Type().FieldByName("_")
	if !ok {
		panic("fillCols: struct must have a metadata field")
	}
	pk := metadata.Tag.Get("pk")
	var rememberPkVal any
	for j := 0; j < val.NumField(); j++ {
		field := val.Field(j)
		fieldType := val.Type().Field(j)
		tagType := fieldType.Tag.Get("type")
		if tagType == "column" {
			col := fieldType.Tag.Get("col")

			colKey := fmt.Sprintf("%s_%s", TableName, col)

			// int64
			if field.Kind() == reflect.Int64 {
				switch v := row[colKey].(type) {
				case int64:
					field.SetInt(v)
				case int: // на случай, если значение int
					field.SetInt(int64(v))
				default:
					field.SetInt(0) // если тип не подходит
				}
				if col == pk {
					rememberPkVal = field.Int()
				}
			}
			// string
			if field.Kind() == reflect.String {
				switch v := row[colKey].(type) {
				case string:
					field.SetString(string(v))
				default:
					field.SetString("") // если тип не подходит
				}
				if col == pk {
					rememberPkVal = field.String()
				}
			}

		}
	}
	return rememberPkVal

}

func (q *Query[T]) smartHydration(model interface{}, parentData []map[string]any, parentPkMap map[string][]any, groupByCol string) map[any]reflect.Value {

	val := reflect.ValueOf(model).Elem()
	valType := val.Type()
	metadata, ok := valType.FieldByName("_")
	if !ok {
		panic("fillCols: struct must have a metadata field")
	}
	TableName := metadata.Tag.Get("table")
	if val.Kind() != reflect.Struct {
		panic("smartHydration: model must be a pointer to a struct")
	}

	structData := map[any]reflect.Value{}
	for _, row := range parentData {
		newStruct := reflect.New(reflect.TypeOf(model).Elem()).Elem()
		pk := q.fillCols(newStruct.Addr().Interface(), row, TableName)
		structData[pk] = newStruct
	}

	for i := 0; i < val.NumField(); i++ {
		relFieldValue := val.Field(i)
		relFieldType := val.Type().Field(i)
		tagType := relFieldType.Tag.Get("type")
		fk := relFieldType.Tag.Get("fk")
		originalKey := relFieldType.Tag.Get("originalkey")
		originalKey = fmt.Sprintf("%s_%s", TableName, originalKey)
		relatedTable := relFieldType.Tag.Get("table")
		relType := relFieldType.Tag.Get("reltype")

		if tagType == "relation" {

			if relType == consts.HasMany && relFieldValue.Kind() == reflect.Slice {
				relStruct := reflect.New(relFieldValue.Type().Elem()).Interface()
				subQuery, subQueryPkMap := NewQuery[T]().SelectFromStruct(relStruct).WhereIn(fk, parentPkMap[originalKey]).LoadMap()
				data := q.smartHydration(relStruct, subQuery, subQueryPkMap, fmt.Sprintf("%s_%s", relatedTable, fk))

				// sliceValue := reflect.MakeSlice(relFieldValue.Type(), 0, len(data))
				// for _, v := range data {
				// 	sliceValue = reflect.Append(sliceValue, v)
				// }
			}

			// if relType == consts.HasOne && relFieldValue.Kind() == reflect.Struct {
			// 	relStruct := reflect.New(relFieldValue.Type()).Interface()
			// 	subQuery, subQueryPkMap := NewQuery[T]().SelectFromStruct(relStruct).WhereIn(fk, parentPkMap[originalKey]).LoadMap()
			// 	q.smartHydration(relStruct, subQuery, subQueryPkMap, groupByCol)
			// 	relFieldValue.Set(reflect.ValueOf(relStruct).Elem())
			// }

			// if relType == consts.BelongsTo && relFieldValue.Kind() == reflect.Struct {
			// 	relStruct := reflect.New(relFieldValue.Type()).Interface()
			// 	subQuery, subQueryPkMap := NewQuery[T]().SelectFromStruct(relStruct).WhereIn(originalKey, parentPkMap[fk]).LoadMap()
			// 	q.smartHydration(relStruct, subQuery, subQueryPkMap, groupByCol)
			// 	relFieldValue.Set(reflect.ValueOf(relStruct).Elem())
			// }
		}
	}

	return structData

}

func (q *Query[T]) Load() []T {
	var model T
	parentData, parentPkMap := q.LoadMap()
	data := q.smartHydration(&model, parentData, parentPkMap)
	return q.reflectSliceToSlice(data)
}
