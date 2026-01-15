package main

import (
	"fmt"
	"reflect"
)

type Post struct {
	Id   int64  `type:"column" col:"posts_id"`
	Name string `type:"column" col:"posts_name"`
}

type Profile struct {
	Id    int64  `type:"column" col:"users_id"`
	Name  string `type:"column" col:"users_name"`
	Posts []Post `type:"relation" table:"posts" reltype:"one-to-many" fk:"users_post_id" rk:"posts_id"`
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

	for i := range val.NumField() {
		value := val.Field(i)
		if value.Kind() == reflect.Int64 {
			//
		}
		fmt.Println(value)
	}

	return nil
}

func main() {
	// ExamplesQuery()
	ExamplesORM()

	ParseStruct(&Profile{Id: 11, Name: "Azgor"})
}
