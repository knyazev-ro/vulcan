package query

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan/orm/model"
)

type Query struct {
	Model         model.Model
	selectExp     string
	whereExp      string
	joinExp       string
	createExp     string
	orderExp      string
	fullStatement string
}

func NewQuery(model model.Model) *Query {
	return &Query{
		Model: model,
	}
}

func (q *Query) Build() *Query {

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

	q.fullStatement += " " + strings.Trim(q.joinExp, " ")
	q.fullStatement = strings.Trim(q.fullStatement, " ")

	q.fullStatement += " " + strings.Trim(q.whereExp, " ")
	q.fullStatement = strings.Trim(q.fullStatement, " ")

	q.fullStatement += " " + strings.Trim(q.orderExp, " ")
	q.fullStatement = strings.Trim(q.fullStatement, " ")

	q.fullStatement += ";"

	return q
}

func (q *Query) SQL() string {
	return q.fullStatement
}

func (q *Query) RawSQL(v string) *Query {
	q.fullStatement = v
	return q
}

func (q *Query) Get() []string {
	return []string{}
}
