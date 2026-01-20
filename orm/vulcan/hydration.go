package vulcan

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/knyazev-ro/vulcan/orm/consts"
)

func (q *Query[T]) reflectSliceToSlice(v []reflect.Value) []T {
	data := make([]T, len(v))
	for i, val := range v {
		data[i] = val.Interface().(T)
	}
	return data
}

func (q *Query[T]) getPk(pk, table string) []string {
	pkArr := []string{}
	for _, k := range strings.Split(strings.Trim(pk, " "), ",") {
		pkArr = append(pkArr, fmt.Sprintf("%s_%s", table, strings.Trim(k, " ")))
	}
	return pkArr
}

func (q *Query[T]) recHydration(data []map[string]any, i reflect.Value) []reflect.Value {
	modelMetadata, ok := i.Type().Elem().FieldByName("_")

	if !ok {
		panic("model metadata was not found")
	}

	pkStr := modelMetadata.Tag.Get("pk")
	tableName := strings.Trim(modelMetadata.Tag.Get("table"), " ")
	pk := q.getPk(pkStr, tableName)

	mapByPk := map[string][]map[string]any{}
	for _, row := range data {
		pKeyVal := []string{}
		for _, pkEl := range pk {
			pKeyVal = append(pKeyVal, fmt.Sprintf("%d", row[pkEl]))
		}
		pKeyValStr := strings.Join(pKeyVal, ",")
		mapByPk[pKeyValStr] = append(mapByPk[pKeyValStr], row)
	}

	structData := []reflect.Value{}

	for _, rowsByPk := range mapByPk {
		newModel := reflect.New(i.Type().Elem()).Elem()
		valModelType := newModel.Type()
		first := rowsByPk[0]
		for i := range newModel.NumField() {
			field := newModel.Field(i)
			fieldType := valModelType.Field(i)
			tagType := fieldType.Tag.Get("type")

			if tagType == "column" {
				col := fieldType.Tag.Get("col")
				colKey := fmt.Sprintf("%s_%s", tableName, col)

				// int64
				if field.Kind() == reflect.Int64 {
					switch v := first[colKey].(type) {
					case int64:
						field.SetInt(v)
					case int: // на случай, если значение int
						field.SetInt(int64(v))
					default:
						field.SetInt(0) // если тип не подходит
					}
				}
				// string
				if field.Kind() == reflect.String {
					switch v := first[colKey].(type) {
					case string:
						field.SetString(string(v))
					default:
						field.SetString("") // если тип не подходит
					}
				}

			}

			if tagType == "relation" {
				relTypeTag := fieldType.Tag.Get("reltype")
				if relTypeTag == consts.HasMany && field.Kind() == reflect.Slice {
					sliceType := field.Type()    // []Some
					elemType := sliceType.Elem() // Some
					childrens := q.recHydration(rowsByPk, reflect.New(elemType))
					newSlice := reflect.MakeSlice(sliceType, 0, len(childrens))
					for _, child := range childrens {
						newSlice = reflect.Append(newSlice, child)
					}
					field.Set(newSlice)
				}

				if relTypeTag == consts.BelongsTo && field.Kind() == reflect.Struct {
					fieldRelationType := field.Type()
					childrens := q.recHydration(rowsByPk, reflect.New(fieldRelationType))
					if len(childrens) >= 1 {
						field.Set(childrens[0])
					}
				}

				if relTypeTag == consts.HasOne && field.Kind() == reflect.Struct {
					fieldRelationType := field.Type()
					childrens := q.recHydration(rowsByPk, reflect.New(fieldRelationType))
					if len(childrens) >= 1 {
						field.Set(childrens[0])
					}
				}
			}
		}
		structData = append(structData, newModel)
	}
	return structData
}

func (q *Query[T]) Hydration(data []map[string]any) []T {
	var m T
	total := q.recHydration(data, reflect.ValueOf(&m))
	return q.reflectSliceToSlice(total)
}
