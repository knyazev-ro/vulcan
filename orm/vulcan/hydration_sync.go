package vulcan

import (
	"context"
	"fmt"
	"reflect"

	"github.com/knyazev-ro/vulcan/orm/consts"
)

// Умная гидрация. Умная, потому что тупая версия реализована была в первом прототипе, использовала LEFT JOIN для отношений
// Умная версия использует WHERE ANY и группировку уже по ним. Хоть вместо 1 запроса будет N, где N - количество отношений в указанной структуре
// Она все равно будет быстрее при малых запросах и феноменально быстрее при больших
func (q *Query[T]) smartHydrationSync(ctx context.Context, model interface{}, parentData []map[string]any, parentPkMap map[string][]any) ([]reflect.Value, error) {

	val := reflect.ValueOf(model).Elem()
	valType := val.Type()
	metadata, ok := valType.FieldByName("_")
	if !ok {
		panic("fillCols: struct must have a metadata field")
	}
	TableName := metadata.Tag.Get("table")
	if val.Kind() != reflect.Struct {
		panic("smartHydration: model must be a pointer to a struct")
	}

	// На этом этапе записываются не-relation поля и создаются сущности родителя, сразу в массив
	structData := []reflect.Value{}
	cachedCols := map[int]string{}
	cachedTags := map[int]map[string]string{}

	for _, row := range parentData {
		newStruct := reflect.New(reflect.TypeOf(model).Elem()).Elem()

		if newStruct.Kind() != reflect.Struct {
			panic("fillCols: i must be a pointer to a struct")
		}

		for j := 0; j < newStruct.NumField(); j++ {
			field := newStruct.Field(j)

			cTags, ex := cachedTags[j]
			tagType := ""
			tableTag := ""
			col := ""
			agg := ""
			if !ex {
				fieldType := newStruct.Type().Field(j)
				tagType = fieldType.Tag.Get("type")
				tableTag = fieldType.Tag.Get("table")
				col = fieldType.Tag.Get("col")
				agg = fieldType.Tag.Get("agg")

				cachedTags[j] = map[string]string{
					"tagType":  tagType,
					"tableTag": tableTag,
					"col":      col,
					"agg":      agg,
				}

			} else {
				tagType = cTags["tagType"]
				tableTag = cTags["tableTag"]
				col = cTags["col"]
				agg = cTags["agg"]
			}

			if tagType == "column" {

				colKey := ""
				if cachedCols[j] == "" {

					if col == "*" {
						col = "all"
					}

					if agg != "" {
						col = fmt.Sprintf(`%s_%s`, col, agg)
					}

					colKey = fmt.Sprintf("%s_%s", TableName, col)
					if tableTag != "" {
						colKey = fmt.Sprintf("%s_%s", tableTag, col)
					}
					cachedCols[j] = colKey
				} else {
					colKey = cachedCols[j]
				}

				// normal primitive values like int string bool etc.
				q.fillWithPrimitive(field, row, colKey)

				// NULL Support!
				if field.Kind() == reflect.Ptr {
					if row[colKey] != nil {
						tmpField := reflect.New(field.Type().Elem())
						q.fillWithPrimitive(tmpField.Elem(), row, colKey)
						field.Set(tmpField)
					} else {
						field.Set(reflect.Zero(field.Type()))
					}
				}

			}
		}

		structData = append(structData, newStruct)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	totalRelations := q.CountRelations(model)
	results := make([]GorutineData, totalRelations)

	for i := 0; i < val.NumField(); i++ {
		relFieldValue := val.Field(i)
		relFieldType := val.Type().Field(i)
		tagType := relFieldType.Tag.Get("type")

		fk := relFieldType.Tag.Get("fk")

		originalKey := relFieldType.Tag.Get("originalkey")
		originalKeyFormatted := fmt.Sprintf("%s_%s", TableName, originalKey)

		relType := relFieldType.Tag.Get("reltype")

		if tagType == "relation" {

			// HasMany
			if relType == consts.HasMany && relFieldValue.Kind() == reflect.Slice {
				// Достаем тип структуры, создаем новую, как образец
				relStructValue := reflect.New(relFieldValue.Type().Elem())
				relStruct := relStructValue.Interface()
				// Выполняем запрос в БД по получению полей.

				whereHasClosure, ok := q.whereHasMap[relFieldType.Name]
				query, err := NewQuery[T]().SelectFromStruct(relStruct).UseConn(q.db)

				if err != nil {
					return nil, err
				}

				if ok {
					whereHasClosure(query)
				}

				subQuery, subQueryPkMap, err := query.WhereAny(fk, parentPkMap[originalKeyFormatted]).LoadMap(ctx)

				if err != nil {
					cancel()
					return nil, err
				}

				// Продолжаем рекурсию
				data, err := q.smartHydrationSync(ctx, relStruct, subQuery, subQueryPkMap)

				if err != nil {
					cancel()
					return nil, err
				}

				// Группируем данные
				dataGrouped := q.groupByKey(data, fk)
				results = append(results, GorutineData{
					dataGrouped:      dataGrouped,
					relFieldTypeName: relFieldType.Name,
					originalKey:      originalKey,
					relType:          consts.HasMany,
				})
			}

			// HasOne
			if relType == consts.HasOne && relFieldValue.Kind() == reflect.Struct {
				relStruct := reflect.New(relFieldValue.Type()).Interface()

				whereHasClosure, ok := q.whereHasMap[relFieldType.Name]
				query, err := NewQuery[T]().SelectFromStruct(relStruct).UseConn(q.db)

				if err != nil {
					return nil, err
				}

				if ok {
					whereHasClosure(query)
				}

				subQuery, subQueryPkMap, err := query.WhereAny(fk, parentPkMap[originalKeyFormatted]).LoadMap(ctx)

				if err != nil {
					cancel()
					return nil, err
				}
				data, err := q.smartHydrationSync(ctx, relStruct, subQuery, subQueryPkMap)

				if err != nil {
					cancel()
					return nil, err
				}
				dataGrouped := q.groupByKey(data, fk)
				results = append(results, GorutineData{
					dataGrouped:      dataGrouped,
					relFieldTypeName: relFieldType.Name,
					originalKey:      originalKey,
					relType:          consts.HasOne,
				})
			}

			// BelongsTo
			if relType == consts.BelongsTo && relFieldValue.Kind() == reflect.Struct {
				relStruct := reflect.New(relFieldValue.Type()).Interface()

				whereHasClosure, ok := q.whereHasMap[relFieldType.Name]
				query, err := NewQuery[T]().SelectFromStruct(relStruct).UseConn(q.db)

				if err != nil {
					return nil, err
				}

				if ok {
					whereHasClosure(query)
				}

				ids := []any{}
				key := fmt.Sprintf("%s_%s", TableName, fk) // o.

				for _, p := range parentData {
					ids = append(ids, p[key])
				}

				subQuery, subQueryPkMap, err := query.WhereAny(originalKey, ids).LoadMap(ctx)

				if err != nil {
					cancel()
					return nil, err
				}

				data, err := q.smartHydrationSync(ctx, relStruct, subQuery, subQueryPkMap)

				if err != nil {
					cancel()
					return nil, err
				}

				dataGrouped := q.groupByKey(data, originalKey)
				results = append(results, GorutineData{
					dataGrouped:      dataGrouped,
					relFieldTypeName: relFieldType.Name,
					fk:               fk,
					relType:          consts.BelongsTo})
			}
		}
	}

	for j := 0; j < totalRelations; j++ {
		r := results[j]
		if r.Error != nil {
			return nil, r.Error
		}
		switch r.relType {
		case consts.HasMany:
			q.placeHasMany(structData, r.dataGrouped, r.originalKey, r.relFieldTypeName)
		case consts.HasOne:
			q.placeHasOne(structData, r.dataGrouped, r.originalKey, r.relFieldTypeName)
		case consts.BelongsTo:
			q.placeBelongsTo(structData, r.dataGrouped, r.fk, r.relFieldTypeName)
		}
	}

	return structData, nil

}
