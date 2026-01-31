package vulcan

import (
	"context"
	"fmt"
	"reflect"

	"github.com/knyazev-ro/vulcan/orm/consts"
	"github.com/knyazev-ro/vulcan/orm/db"
)

// Модель можно загрузить по схеме структуры в виде плоского MAP.
func (q *Query[T]) LoadMap(ctx context.Context) ([]map[string]any, map[string][]any, error) {
	mapData := []map[string]any{}
	pkMap := map[string][]any{}
	q.Build()
	// fmt.Println(q.SQL())
	db := db.DB
	rows, err := db.QueryContext(ctx, q.fullStatement, q.Bindings...)
	if err != nil {
		return mapData, pkMap, err
	}

	defer rows.Close()

	cols, _ := rows.Columns()

	table := q.Model.TableName
	pk := q.Model.Pk

	pkCols := q.getPk(pk, table)

	n := len(cols)
	colValues := make([]any, n)
	colPtrs := make([]any, n)
	for i := range colValues {
		colPtrs[i] = &colValues[i]
	}

	for rows.Next() {
		if err := rows.Scan(colPtrs...); err != nil {
			panic(err)
		}

		colsMap := make(map[string]any, n)
		for i, col := range cols {
			colsMap[col] = colValues[i]
		}

		for _, pkCol := range pkCols {
			_, e := colsMap[pkCol]
			if !e {
				continue
			}
			pkMap[pkCol] = append(pkMap[pkCol], colsMap[pkCol])
		}

		mapData = append(mapData, colsMap)
	}
	return mapData, pkMap, nil
}

// Заполнение обычных полей, игнорирует заполнение полей-отношений

// Группирует данные по ключу. Полезен для группировки по внешнему ключу, либо по собственному ключу сущности
// Принимает ключ сущности по которому производится группировка и массив этих сущностей для операции
func (q *Query[T]) groupByKey(data []reflect.Value, key string) map[any][]reflect.Value {
	grouped := map[any][]reflect.Value{}
	if len(data) <= 0 {
		return grouped
	}
	example := data[0]
	fieldName := ""
	for i := 0; i < example.NumField(); i++ {
		fieldType := example.Type().Field(i)
		col := fieldType.Tag.Get("col")

		if col != key {
			continue
		}
		fieldName = fieldType.Name
	}

	for _, d := range data {
		f := d.FieldByName(fieldName)
		grouped[f.Int()] = append(grouped[f.Int()], d)
	}
	return grouped
}

// Вставляет массив в поле массив для отношений Имеет Много (HasMany)
func (q *Query[T]) placeHasMany(parent []reflect.Value, grouped map[any][]reflect.Value, originalkey string, fieldName string) {
	example := parent[0]
	originalKeyFieldName := ""
	for i := 0; i < example.NumField(); i++ {
		fieldType := example.Type().Field(i)
		col := fieldType.Tag.Get("col")
		if col != originalkey {
			continue
		}
		originalKeyFieldName = fieldType.Name
	}

	for _, p := range parent {
		id := p.FieldByName(originalKeyFieldName).Int()
		relatedArr := grouped[id]
		if len(relatedArr) <= 0 {
			continue
		}

		neededField := p.FieldByName(fieldName)
		relatedType := neededField.Type()

		newSlice := reflect.MakeSlice(relatedType, 0, len(relatedArr))
		for _, child := range relatedArr {
			newSlice = reflect.Append(newSlice, child)
		}

		neededField.Set(newSlice)
	}
}

// Вставляет поле-структуру, реализует логику для HasOne
func (q *Query[T]) placeHasOne(parent []reflect.Value, grouped map[any][]reflect.Value, originalkey string, fieldName string) {
	example := parent[0]
	originalKeyFieldName := ""
	for i := 0; i < example.NumField(); i++ {
		fieldType := example.Type().Field(i)
		col := fieldType.Tag.Get("col")
		if col != originalkey {
			continue
		}
		originalKeyFieldName = fieldType.Name
	}

	for _, p := range parent {
		id := p.FieldByName(originalKeyFieldName).Int()
		relatedArr := grouped[id]
		if len(relatedArr) <= 0 {
			continue
		}

		neededField := p.FieldByName(fieldName)
		neededField.Set(relatedArr[0])
	}
}

// Вставляет поле-структуру, реализует логику для BelongsTo
func (q *Query[T]) placeBelongsTo(parent []reflect.Value, grouped map[any][]reflect.Value, fk string, fieldName string) {
	// Находим у сущности поле связывающее его с родителем
	example := parent[0]
	fkFieldName := ""
	for i := 0; i < example.NumField(); i++ {
		fieldType := example.Type().Field(i)
		col := fieldType.Tag.Get("col")
		if col != fk {
			continue
		}
		fkFieldName = fieldType.Name
	}
	// Сохраняем записи. Проходимся по ранее сгруппированным по ключу полям (id того, кто Владеет)
	// Родитель, у которого совпадает поле внешний ключ с id того-кто-владеет записывает у себя эту сущность
	for _, p := range parent {
		fkId := p.FieldByName(fkFieldName).Int()
		relatedArr := grouped[fkId]
		if len(relatedArr) <= 0 {
			continue
		}

		neededField := p.FieldByName(fieldName)
		neededField.Set(relatedArr[0])

	}
}

func (q *Query[T]) CountRelations(i interface{}) int {
	val := reflect.ValueOf(i).Elem()
	if val.Kind() != reflect.Struct {
		panic("must be a struct!")
	}

	count := 0
	for i := 0; i < val.NumField(); i++ {
		tagType := val.Type().Field(i).Tag.Get("type")
		if tagType == "relation" {
			count++
		}
	}
	return count
}

type GorutineData struct {
	dataGrouped      map[any][]reflect.Value
	relFieldTypeName string
	originalKey      string
	relType          string
	fk               string
	Error            error
}

func (q *Query[T]) fillWithPrimitive(field reflect.Value, row map[string]any, colKey string) reflect.Value {
	if field.Kind() == reflect.Int64 {
		switch v := row[colKey].(type) {
		case int64:
			field.SetInt(v)
		case int: // на случай, если значение int
			field.SetInt(int64(v))
		default:
			field.SetInt(0) // если тип не подходит
		}
		return field

	}
	// string
	if field.Kind() == reflect.String {
		switch v := row[colKey].(type) {
		case string:
			field.SetString(string(v))
		default:
			field.SetString("") // если тип не подходит
		}
		return field
	}

	// bool
	if field.Kind() == reflect.Bool {
		switch v := row[colKey].(type) {
		case bool:
			field.SetBool(bool(v))
		default:
			field.SetBool(false) // если тип не подходит
		}
		return field

	}
	return field
}

// Умная гидрация. Умная, потому что тупая версия реализована была в первом прототипе, использовала LEFT JOIN для отношений
// Умная версия использует WHERE ANY и группировку уже по ним. Хоть вместо 1 запроса будет N, где N - количество отношений в указанной структуре
// Она все равно будет быстрее при малых запросах и феноменально быстрее при больших
func (q *Query[T]) smartHydration(ctx context.Context, model interface{}, parentData []map[string]any, parentPkMap map[string][]any) ([]reflect.Value, error) {

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
	results := make(chan GorutineData, totalRelations)

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
				query := NewQuery[T]().SelectFromStruct(relStruct)
				if ok {
					whereHasClosure(query)
				}

				go func() {
					select {
					case <-ctx.Done():
						results <- GorutineData{Error: ctx.Err()}
						return
					case db.GlobalLimit <- struct{}{}:
						// Wait...
					}

					defer func() { <-db.GlobalLimit }()

					subQuery, subQueryPkMap, err := query.WhereAny(fk, parentPkMap[originalKeyFormatted]).LoadMap(ctx)

					if err != nil {
						cancel()
						results <- GorutineData{Error: err}
						return
					}

					// Продолжаем рекурсию
					data, err := q.smartHydration(ctx, relStruct, subQuery, subQueryPkMap)

					if err != nil {
						cancel()
						results <- GorutineData{Error: err}
						return
					}
					// Группируем данные
					dataGrouped := q.groupByKey(data, fk)
					results <- GorutineData{
						dataGrouped:      dataGrouped,
						relFieldTypeName: relFieldType.Name,
						originalKey:      originalKey,
						relType:          consts.HasMany,
					}
				}()
				// Выполняем вставку по полям родителя в соответсвующий relation в памяти
				// q.placeHasMany(structData, dataGrouped, originalKey, relFieldType.Name)
			}

			// HasOne
			if relType == consts.HasOne && relFieldValue.Kind() == reflect.Struct {
				relStruct := reflect.New(relFieldValue.Type()).Interface()

				whereHasClosure, ok := q.whereHasMap[relFieldType.Name]
				query := NewQuery[T]().SelectFromStruct(relStruct)
				if ok {
					whereHasClosure(query)
				}

				go func() {
					select {
					case <-ctx.Done():
						results <- GorutineData{Error: ctx.Err()}
						return
					case db.GlobalLimit <- struct{}{}:
						// Wait...
					}

					defer func() { <-db.GlobalLimit }()

					subQuery, subQueryPkMap, err := query.WhereAny(fk, parentPkMap[originalKeyFormatted]).LoadMap(ctx)

					if err != nil {
						cancel()
						results <- GorutineData{Error: err}
						return
					}
					data, err := q.smartHydration(ctx, relStruct, subQuery, subQueryPkMap)

					if err != nil {
						cancel()
						results <- GorutineData{Error: err}
						return
					}
					dataGrouped := q.groupByKey(data, fk)
					results <- GorutineData{
						dataGrouped:      dataGrouped,
						relFieldTypeName: relFieldType.Name,
						originalKey:      originalKey,
						relType:          consts.HasOne,
					}
				}()

				// q.placeHasOne(structData, dataGrouped, originalKey, relFieldType.Name)
			}

			// BelongsTo
			if relType == consts.BelongsTo && relFieldValue.Kind() == reflect.Struct {
				relStruct := reflect.New(relFieldValue.Type()).Interface()

				whereHasClosure, ok := q.whereHasMap[relFieldType.Name]
				query := NewQuery[T]().SelectFromStruct(relStruct)
				if ok {
					whereHasClosure(query)
				}

				ids := []any{}
				key := fmt.Sprintf("%s_%s", TableName, fk) // o.

				for _, p := range parentData {
					ids = append(ids, p[key])
				}

				go func() {
					select {
					case <-ctx.Done():
						results <- GorutineData{Error: ctx.Err()}
						return
					case db.GlobalLimit <- struct{}{}:
						// Wait...
					}

					defer func() { <-db.GlobalLimit }()

					subQuery, subQueryPkMap, err := query.WhereAny(originalKey, ids).LoadMap(ctx)

					if err != nil {
						cancel()
						results <- GorutineData{Error: err}
						return
					}
					data, err := q.smartHydration(ctx, relStruct, subQuery, subQueryPkMap)
					if err != nil {
						cancel()
						results <- GorutineData{Error: err}
						return
					}
					dataGrouped := q.groupByKey(data, originalKey)
					results <- GorutineData{
						dataGrouped:      dataGrouped,
						relFieldTypeName: relFieldType.Name,
						fk:               fk,
						relType:          consts.BelongsTo}
				}()
				// q.placeBelongsTo(structData, dataGrouped, fk, relFieldType.Name)
			}
		}
	}

	for j := 0; j < totalRelations; j++ {
		r := <-results

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

func (q *Query[T]) Load(ctx context.Context) ([]T, error) {
	var model T
	parentData, parentPkMap, err := q.LoadMap(ctx)
	if err != nil {
		return []T{}, err
	}
	data, err := q.smartHydration(ctx, &model, parentData, parentPkMap)
	if err != nil {
		return nil, err
	}
	return q.reflectSliceToSlice(data), nil
}

func (q *Query[T]) With(field string, closure func(*Query[T])) *Query[T] {
	if q.whereHasMap == nil {
		q.whereHasMap = make(map[string]func(*Query[T]))
	}
	q.whereHasMap[field] = closure
	return q
}
