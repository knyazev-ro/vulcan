package query

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/knyazev-ro/vulcan/config"
	"github.com/knyazev-ro/vulcan/orm/model"
)

const OneToMany string = "one-to-many"

func (q *Query[T]) recGenerateCols(i interface{}, cols []string) []string {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		panic("Must be a struct")
	}
	pk := ""
	TableName := ""
	for i := range val.NumField() {
		field := val.Field(i)
		valueType := val.Type().Field(i)
		typeTag := valueType.Tag.Get("type")

		if typeTag == "metadata" {
			TableName = valueType.Tag.Get("table")
			pk = valueType.Tag.Get("pk")
		}

		if typeTag == "column" {
			colTag := valueType.Tag.Get("col")
			cols = append(cols, fmt.Sprintf(`"%s"."%s" AS %s_%s`, TableName, colTag, TableName, colTag))
		}

		if typeTag == "relation" {
			relTypeTag := valueType.Tag.Get("reltype")
			tableTag := valueType.Tag.Get("table")
			fkTag := valueType.Tag.Get("fk")
			// one to many
			if relTypeTag == OneToMany && field.Kind() == reflect.Slice {
				cols = q.recGenerateCols(reflect.New(field.Type().Elem()).Interface(), cols)
				q.LeftJoin(tableTag, func(jc *Join) {
					jc.On(fmt.Sprintf(`%s.%s`, tableTag, fkTag), "=", fmt.Sprintf(`%s.%s`, TableName, pk))
				})
			}
		}
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

func (q *Query[T]) getColsPtr(i interface{}) []any {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		panic("Must be a struct")
	}

	colsPtr := []any{}

	for i := range val.NumField() {
		value := val.Field(i)
		valueType := val.Type().Field(i)
		typeTag := valueType.Tag.Get("type")
		if typeTag == "column" {
			colsPtr = append(colsPtr, value.Addr().Interface())
		}
		if typeTag == "relation" {
			tableTag := valueType.Tag.Get("table")
			relTypeTag := valueType.Tag.Get("reltype")
			fkTag := valueType.Tag.Get("fk")
			fmt.Println(typeTag, tableTag, relTypeTag, fkTag)
			// logic
		}
	}

	return colsPtr
}

func (q *Query[T]) Get() []T {
	config := config.GetConfig()
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s", config.Driver, config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("pgx", dsn) // pgx через database/sql
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		panic(err)
	}
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
		fmt.Println(colsMap)
		mapData = append(mapData, colsMap)
	}

	return q.HydrationOneToMany(mapData)
}
