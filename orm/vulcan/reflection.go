package vulcan

import (
	"fmt"
	"reflect"

	"github.com/knyazev-ro/vulcan/orm/model"
)

func (q *Query[T]) generateCols(i interface{}, cols []string) []string {
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
	TableName := metadata.Tag.Get("table")

	for i := range val.NumField() {
		valueType := val.Type().Field(i)
		typeTag := valueType.Tag.Get("type")

		if typeTag == "column" {
			colTag := valueType.Tag.Get("col")
			tableTag := valueType.Tag.Get("table")
			if tableTag == "" {
				cols = append(cols, fmt.Sprintf(`"%s"."%s" AS %s_%s`, TableName, colTag, TableName, colTag))
			} else {
				cols = append(cols, fmt.Sprintf(`"%s"."%s" AS %s_%s`, tableTag, colTag, tableTag, colTag))
			}
		}
	}
	return cols
}

func (q *Query[T]) MSelect(i interface{}) *Query[T] {
	cols := q.generateCols(i, []string{})
	metadata, ok := reflect.TypeOf(i).Elem().FieldByName("_")
	if !ok {
		panic("metadata is not found")
	}
	q.Model = model.Model{
		TableName: metadata.Tag.Get("table"),
		Pk:        metadata.Tag.Get("pk"),
	}
	if len(cols) > 0 {
		q.selectRaw(cols)
	}
	return q
}

// func (q *Query[T]) DeprectedGet() []T {
// 	db := db.DB
// 	println(q.SQL())
// 	rows, err := db.Query(q.fullStatement, q.Bindings...)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer rows.Close()

// 	cols, _ := rows.Columns()
// 	fmt.Println(rows.Columns())
// 	mapData := []map[string]any{}
// 	for rows.Next() {

// 		colValues := make([]any, len(cols))
// 		colPtrs := make([]any, len(cols))
// 		colsMap := map[string]any{}

// 		for i := range colValues {
// 			colPtrs[i] = &colValues[i]
// 		}
// 		if err := rows.Scan(colPtrs...); err != nil {
// 			panic(err)
// 		}

// 		for i, col := range cols {
// 			colsMap[col] = colValues[i]
// 		}
// 		mapData = append(mapData, colsMap)
// 	}

// 	return q.Hydration(mapData)
// }
