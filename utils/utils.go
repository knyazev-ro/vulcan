package utils

import (
	"fmt"
	"reflect"
)

func GetFieldInfo(s interface{}) []map[string]string {
	t := reflect.TypeOf(s)
	var result []map[string]string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		info := map[string]string{
			"name": field.Name,
			"type": field.Type.String(),
			"tag":  string(field.Tag),
		}
		result = append(result, info)
	}

	return result
}

func ColsSafe(cols []string) []string {
	colsSafe := []string{}
	for _, col := range cols {
		colsSafe = append(colsSafe, fmt.Sprintf(`"%s"`, col)) // psql
	}
	return colsSafe
}
