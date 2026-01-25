package vulcan

import (
	"fmt"
	"reflect"

	"github.com/knyazev-ro/vulcan/orm/consts"
	"github.com/knyazev-ro/vulcan/orm/db"
)

func MapByKey(data []map[string]any, key string) map[any][]map[string]any {
	result := make(map[any][]map[string]any)
	for _, item := range data {
		if val, exists := item[key]; exists {
			result[val] = append(result[val], item)
		}
	}
	return result
}

// Модель можно загрузить по схеме структуры в виде плоского MAP.
func (q *Query[T]) LoadMap() ([]map[string]any, map[string][]any) {
	q.Build()
	fmt.Println(q.SQL())
	db := db.DB
	rows, err := db.Query(q.fullStatement, q.Bindings...)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	cols, _ := rows.Columns()

	table := q.Model.TableName
	pk := q.Model.Pk

	pkCols := q.getPk(pk, table)
	mapData := []map[string]any{}
	pkMap := map[string][]any{}
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

		for _, pkCol := range pkCols {
			pkMap[pkCol] = append(pkMap[pkCol], colsMap[pkCol])
		}

		mapData = append(mapData, colsMap)
	}
	return mapData, pkMap
}

// Заполнение обычных полей, игнорирует заполнение полей-отношений
func (q *Query[T]) fillCols(i interface{}, row map[string]any, TableName string) any {
	val := reflect.ValueOf(i).Elem()
	if val.Kind() != reflect.Struct {
		panic("fillCols: i must be a pointer to a struct")
	}

	metadata, ok := val.Type().FieldByName("_")
	if !ok {
		panic("fillCols: struct must have a metadata field")
	}

	var rememberPkVal any
	pk := metadata.Tag.Get("pk")
	for j := 0; j < val.NumField(); j++ {
		field := val.Field(j)
		fieldType := val.Type().Field(j)
		tagType := fieldType.Tag.Get("type")
		if tagType == "column" {
			col := fieldType.Tag.Get("col")

			colKey := fmt.Sprintf("%s_%s", TableName, col)

			// int64
			if field.Kind() == reflect.Int64 {
				switch v := row[colKey].(type) {
				case int64:
					field.SetInt(v)
				case int: // на случай, если значение int
					field.SetInt(int64(v))
				default:
					field.SetInt(0) // если тип не подходит
				}
				if col == pk {
					rememberPkVal = field.Int()
				}
			}
			// string
			if field.Kind() == reflect.String {
				switch v := row[colKey].(type) {
				case string:
					field.SetString(string(v))
				default:
					field.SetString("") // если тип не подходит
				}
				if col == pk {
					rememberPkVal = field.String()
				}
			}

		}
	}
	return rememberPkVal

}

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

// Умная гидрация. Умная, потому что тупая версия реализована была в первом прототипе, использовала LEFT JOIN для отношений
// Умная версия использует WHERE ANY и группировку уже по ним. Хоть вместо 1 запроса будет N, где N - количество отношений в указанной структуре
// Она все равно будет быстрее при малых запросах и феноменально быстрее при больших
func (q *Query[T]) smartHydration(model interface{}, parentData []map[string]any, parentPkMap map[string][]any) []reflect.Value {

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
	for _, row := range parentData {
		newStruct := reflect.New(reflect.TypeOf(model).Elem()).Elem()
		q.fillCols(newStruct.Addr().Interface(), row, TableName)
		structData = append(structData, newStruct)
	}

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

				subQuery, subQueryPkMap := query.WhereAny(fk, parentPkMap[originalKeyFormatted]).LoadMap()
				// Продолжаем рекурсию
				data := q.smartHydration(relStruct, subQuery, subQueryPkMap)
				// Группируем данные
				dataGrouped := q.groupByKey(data, fk)
				// Выполняем вставку по полям родителя в соответсвующий relation в памяти
				q.placeHasMany(structData, dataGrouped, originalKey, relFieldType.Name)
			}

			// HasOne
			if relType == consts.HasOne && relFieldValue.Kind() == reflect.Struct {
				relStruct := reflect.New(relFieldValue.Type()).Interface()

				whereHasClosure, ok := q.whereHasMap[relFieldType.Name]
				query := NewQuery[T]().SelectFromStruct(relStruct)
				if ok {
					whereHasClosure(query)
				}

				subQuery, subQueryPkMap := query.WhereAny(fk, parentPkMap[originalKeyFormatted]).LoadMap()
				data := q.smartHydration(relStruct, subQuery, subQueryPkMap)
				dataGrouped := q.groupByKey(data, fk)
				q.placeHasOne(structData, dataGrouped, originalKey, relFieldType.Name)
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
				for _, p := range parentData {
					key := fmt.Sprintf("%s_%s", TableName, fk)
					ids = append(ids, p[key])
				}
				subQuery, subQueryPkMap := query.WhereAny(originalKey, ids).LoadMap()
				data := q.smartHydration(relStruct, subQuery, subQueryPkMap)
				dataGrouped := q.groupByKey(data, originalKey)
				q.placeBelongsTo(structData, dataGrouped, fk, relFieldType.Name)
			}
		}
	}

	return structData

}

func (q *Query[T]) Load() []T {
	var model T
	parentData, parentPkMap := q.LoadMap()
	data := q.smartHydration(&model, parentData, parentPkMap)
	return q.reflectSliceToSlice(data)
}

func (q *Query[T]) With(field string, closure func(*Query[T])) *Query[T] {
	if q.whereHasMap == nil {
		q.whereHasMap = make(map[string]func(*Query[T]))
	}
	q.whereHasMap[field] = closure
	return q
}
