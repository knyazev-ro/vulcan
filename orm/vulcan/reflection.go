package vulcan

import (
	"fmt"
	"reflect"

	"github.com/knyazev-ro/vulcan/orm/db"
	"github.com/knyazev-ro/vulcan/orm/model"
)

func (q *Query[T]) recGenerateCols(i interface{}, cols []string) []string {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		panic("Must be a struct")
	}
	metadata, ok := val.Type().FieldByName("_")
	if !ok {
		panic("metadata is not found")
	}
	// pk := metadata.Tag.Get("pk")
	TableName := metadata.Tag.Get("table")

	for i := range val.NumField() {
		// field := val.Field(i)
		valueType := val.Type().Field(i)
		typeTag := valueType.Tag.Get("type")

		if typeTag == "column" {
			colTag := valueType.Tag.Get("col")
			cols = append(cols, fmt.Sprintf(`"%s"."%s" AS %s_%s`, TableName, colTag, TableName, colTag))
		}

		// if typeTag == "relation" {
		// 	relTypeTag := valueType.Tag.Get("reltype")
		// 	tableTag := valueType.Tag.Get("table")
		// 	fkTag := valueType.Tag.Get("fk")
		// 	originalKey := valueType.Tag.Get("originalkey")
		// 	// one to many
		// 	if relTypeTag == consts.HasMany && field.Kind() == reflect.Slice {
		// 		q.LeftJoin(tableTag, func(jc *Join) {
		// 			jc.On(fmt.Sprintf(`%s.%s`, tableTag, fkTag), "=", fmt.Sprintf(`%s.%s`, TableName, pk))
		// 		})
		// 		cols = q.recGenerateCols(reflect.New(field.Type().Elem()).Interface(), cols)
		// 	}
		// 	if relTypeTag == consts.BelongsTo && field.Kind() == reflect.Struct {
		// 		q.LeftJoin(tableTag, func(jc *Join) {
		// 			jc.On(fmt.Sprintf(`%s.%s`, tableTag, originalKey), "=", fmt.Sprintf(`%s.%s`, TableName, fkTag))
		// 		})
		// 		cols = q.recGenerateCols(reflect.New(field.Type()).Interface(), cols)
		// 	}

		// 	if relTypeTag == consts.HasOne && field.Kind() == reflect.Struct {
		// 		q.LeftJoin(tableTag, func(jc *Join) {
		// 			jc.On(fmt.Sprintf(`%s.%s`, tableTag, fkTag), "=", fmt.Sprintf(`%s.%s`, TableName, pk))
		// 		})
		// 		cols = q.recGenerateCols(reflect.New(field.Type()).Interface(), cols)
		// 	}

		// }
	}
	return cols
}

func (q *Query[T]) MSelect(i interface{}) *Query[T] {
	cols := q.recGenerateCols(i, []string{})
	TName, ok := reflect.TypeOf(i).Elem().FieldByName("_")
	if !ok {
		panic("metadata is not found")
	}
	q.Model = model.Model{
		TableName: TName.Tag.Get("table"),
	}
	if len(cols) > 0 {
		q.selectRaw(cols)
	}
	return q
}

func (q *Query[T]) Get() []T {
	db := db.DB
	println(q.SQL())
	rows, err := db.Query(q.fullStatement, q.Bindings...)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	cols, _ := rows.Columns()
	fmt.Println(rows.Columns())
	mapData := []map[string]any{}
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
		mapData = append(mapData, colsMap)
	}

	return q.Hydration(mapData)
}
