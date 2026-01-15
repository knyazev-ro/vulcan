package models

import "reflect"

type Post struct {
	Id   int64  `db:"posts_id"`
	Name string `db:"posts_name"`
}

type User struct {
	Id    int64  `db:"users_id"`
	Name  string `db:"users_name"`
	Posts []Post `fk:"users_post_id"`
}

func ParseStruct(i interface{}) {
	val := reflect.ValueOf(i)
	// t := reflect.TypeOf(i)
	println(val.Type())
}
