package vulcan

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // драйвер для database/sql
	"github.com/knyazev-ro/vulcan/orm/db"
	"github.com/knyazev-ro/vulcan/orm/model"
)

type Query[T any] struct {
	returnType    T
	Model         model.Model
	Bindings      []any
	selectExp     string
	whereExp      string
	joinExp       string
	createExp     string
	orderExp      string
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
	q.MSelect(&i)
	return q
}

func (q *Query[T]) SelectFromStruct(i interface{}) *Query[T] {
	cols := q.generateCols(i, &GenerateColsOptions{useAggs: true})
	if len(cols) > 0 {
		q.selectRaw(cols)
	}
	metadata, ok := reflect.TypeOf(i).Elem().FieldByName("_")
	if !ok {
		panic("metadata is not found")
	}
	q.Model = model.Model{
		TableName: metadata.Tag.Get("table"),
		Pk:        metadata.Tag.Get("pk"),
	}
	return q
}

func (q *Query[T]) Build() *Query[T] {

	if q.fullStatement != "" {
		return q
	}

	selectStr := fmt.Sprintf("SELECT *")

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

	if q.orderExp != "" {
		q.fullStatement += " " + strings.Trim(q.orderExp, " ")
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
	// println("SQL: ", q.fullStatement)
	return q.fullStatement
}

// опасно
func (q *Query[T]) RawSQL(v string) *Query[T] {
	q.fullStatement = v
	return q
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
