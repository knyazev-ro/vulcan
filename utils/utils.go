package utils

import (
	"fmt"
	"reflect"
	"strings"
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

func SeparateParts(side string) string {
	sideParts := strings.Split(side, ".")
	secPart := ""
	if len(sideParts) == 2 {
		secPart = fmt.Sprintf(`."%s"`, sideParts[1])
	}
	return fmt.Sprintf(`"%s"%s`, sideParts[0], secPart)
}

func ColsSafe(cols []string) []string {
	colsSafe := []string{}
	for _, col := range cols {
		colNormialize := SeparateParts(col)
		colsSafe = append(colsSafe, fmt.Sprintf(`%s`, colNormialize)) // psql
	}
	return colsSafe
}
