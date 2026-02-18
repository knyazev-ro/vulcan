package vulcan

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // драйвер для database/sql
	"github.com/knyazev-ro/vulcan-orm/orm/db"
	"github.com/knyazev-ro/vulcan-orm/orm/model"
)

type Query[T any] struct {
	Model         model.Model
	Bindings      []any
	selectExp     string
	whereExp      string
	joinExp       string
	createExp     string
	orderExp      []string
	fromExp       string
	usingExp      string
	limitExp      string
	offsetExp     string
	groupByExp    string
	fullStatement string
	whereHasMap   map[string]func(*Query[T])

	db DBConnection
}

func NewQuery[T any]() *Query[T] {
	q := &Query[T]{}
	q.db = db.DB
	var i T
	q.SelectFromStruct(&i)
	return q
}

func (q *Query[T]) getFieldForCorrespondedColTags(i interface{}, cols []string) []string {
	val := reflect.ValueOf(i).Elem()
	structFieldNames := []string{}
	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i)
		col := fieldType.Tag.Get("col")
		if col == "" {
			continue
		}

		for _, colElem := range cols {
			if col == colElem {
				structFieldNames = append(structFieldNames, fieldType.Name)
				break
			}
		}
	}
	return structFieldNames
}

func (q *Query[T]) SelectFromStruct(i interface{}) *Query[T] {
	cols := q.generateCols(i, &GenerateColsOptions{useAggs: true})
	metadata, ok := reflect.TypeOf(i).Elem().FieldByName("_")
	if !ok {
		panic("metadata is not found")
	}
	pk := metadata.Tag.Get("pk")
	pks := strings.Split(pk, ",")
	pksInStruct := q.getFieldForCorrespondedColTags(i, pks)
	q.Model = model.Model{
		TableName:   metadata.Tag.Get("table"),
		Pk:          pk,
		Pks:         pks,
		PksInStruct: pksInStruct,
	}

	if len(cols) > 0 {
		q.selectRaw(cols)
	}
	return q
}

func (q *Query[T]) Build() *Query[T] {

	if q.fullStatement != "" {
		return q
	}

	selectStr := "SELECT *"

	if q.selectExp != "" {
		selectStr = q.selectExp
	}

	from := fmt.Sprintf("FROM %s", q.Model.TableName)

	q.fullStatement = fmt.Sprintf("%s %s", strings.Trim(selectStr, " "), from)
	q.fullStatement = strings.Trim(q.fullStatement, " ")
	q.appendExpressions()
	q.fullStatement += ";"

	q.fillBindingsPSQL()

	return q
}

func (q *Query[T]) appendExpressions() {

	if q.fromExp != "" {
		q.fullStatement += " " + strings.Trim(q.fromExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.usingExp != "" {
		q.fullStatement += " USING " + strings.Trim(q.usingExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.joinExp != "" {
		q.fullStatement += " " + strings.Trim(q.joinExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.whereExp != "" {
		q.fullStatement += " WHERE " + strings.Trim(q.whereExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.groupByExp != "" {
		q.fullStatement += " " + strings.Trim(q.groupByExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if len(q.orderExp) > 0 {
		q.fullStatement += " ORDER BY " + strings.Trim(strings.Join(q.orderExp, ", "), " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.limitExp != "" {
		q.fullStatement += " " + strings.Trim(q.limitExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.offsetExp != "" {
		q.fullStatement += " " + strings.Trim(q.offsetExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}
}

func (q *Query[T]) SQL() string {
	return q.fullStatement
}

func (q *Query[T]) fillBindingsPSQL() {
	var b strings.Builder
	b.Grow(len(q.fullStatement) + 16)

	idx := 1
	for i := 0; i < len(q.fullStatement); i++ {
		if q.fullStatement[i] == '?' {
			b.WriteString(fmt.Sprintf("$%d", idx))
			idx++
		} else {
			b.WriteByte(q.fullStatement[i])
		}
	}

	q.fullStatement = b.String()
}

func (q *Query[T]) Clone() *Query[T] {
	return &Query[T]{
		Model:         q.Model,
		Bindings:      q.Bindings,
		selectExp:     q.selectExp,
		whereExp:      q.whereExp,
		joinExp:       q.joinExp,
		createExp:     q.createExp,
		orderExp:      q.orderExp,
		fromExp:       q.fromExp,
		usingExp:      q.usingExp,
		limitExp:      q.limitExp,
		offsetExp:     q.offsetExp,
		groupByExp:    q.groupByExp,
		fullStatement: q.fullStatement,
		whereHasMap:   q.whereHasMap,
		db:            q.db,
	}
}
