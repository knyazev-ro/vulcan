package query

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/knyazev-ro/vulcan/config"
	"github.com/knyazev-ro/vulcan/orm/model"
)

func (q *Query[T]) MSelect(i interface{}) *Query[T] {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		panic("Must be a struct")
	}
	cols := []string{}
	TableName := ""

	for i := range val.NumField() {
		valueType := val.Type().Field(i)
		typeTag := valueType.Tag.Get("type")

		if typeTag == "metadata" {
			TableName = valueType.Tag.Get("table")
		}

		if typeTag == "column" {
			colTag := valueType.Tag.Get("col")
			cols = append(cols, fmt.Sprintf(`"%s"."%s" AS %s_%s`, TableName, colTag, TableName, colTag))
		}
	}
	q.Model = model.Model{
		TableName: TableName,
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
	}

	return colsPtr
}

func (q *Query[T]) Get() []T {
	var data []T
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

	fmt.Println(rows.Columns())
	for rows.Next() {
		var m T
		colsPtr := q.getColsPtr(&m)
		if err := rows.Scan(colsPtr...); err != nil {
			panic(err)
		}
		data = append(data, m)
	}
	return data
}
