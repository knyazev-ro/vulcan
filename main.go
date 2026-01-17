package main

import (
	"fmt"
	"reflect"
)

type Post struct {
	TableName string
	Id        int64  `type:"column" col:"id"`
	Name      string `type:"column" col:"name"`
}

type Profile struct {
	TableName string
	Id        int64  `type:"column" col:"id"`
	Name      string `type:"column" col:"name"`
	Posts     []Post `type:"relation" table:"posts" reltype:"one-to-many" fk:"user_id"`
}

type InvalidStructError struct {
	message string
}

func (e *InvalidStructError) Error() string {
	return e.message
}

func ParseStruct(i interface{}) error {
	val := reflect.ValueOf(i)

	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		return &InvalidStructError{message: "Must be a struct"}
	}
	cols := []string{}

	for i := range val.NumField() {
		value := val.Field(i)
		valueType := val.Type().Field(i)
		structFieldName := valueType.Name

		typeTag := valueType.Tag.Get("type")

		if structFieldName == "TableName" {
			fmt.Println("TableName is", value)
		}

		if typeTag == "column" {
			colTag := valueType.Tag.Get("col")
			fmt.Println(typeTag, colTag)
			cols = append(cols, colTag)
			// logic
		}

		if typeTag == "relation" {
			tableTag := valueType.Tag.Get("table")
			relTypeTag := valueType.Tag.Get("reltype")
			fkTag := valueType.Tag.Get("fk")
			fmt.Println(typeTag, tableTag, relTypeTag, fkTag)
			// logic
		}

	}

	return nil
}

func main() {
	// ExamplesQuery()
	ExamplesORM()
}
