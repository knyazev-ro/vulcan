package query

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan/utils"
)

func (q *Query[T]) Select(cols []string) *Query[T] {
	colsSafe := utils.ColsSafe(cols)
	colsStr := strings.Join(colsSafe, ", ")
	selectStatement := fmt.Sprintf("SELECT %s", colsStr)
	q.selectExp = selectStatement
	return q
}

func (q *Query[T]) selectRaw(cols []string) *Query[T] {
	colsStr := strings.Join(cols, ", ")
	selectStatement := fmt.Sprintf("SELECT %s", colsStr)
	q.selectExp = selectStatement
	return q
}

func (q *Query[T]) From(table string) *Query[T] {
	q.fromExp = fmt.Sprintf(`FROM "%s"`, table)
	return q
}

func (q *Query[T]) Using(table string) *Query[T] {
	q.usingExp = fmt.Sprintf(`Using "%s"`, table)
	return q
}

func (q *Query[T]) OnStatment(left string, expr string, right string, clay bool) *Query[T] {
	whereStr := ""
	if q.whereExp == "" {
		whereStr = "WHERE "
	}

	if len(q.whereExp) == 1 && q.whereExp == "(" {
		q.whereExp = "WHERE ("
	}

	statement := fmt.Sprintf(`%s %s %s %s`, whereStr, utils.SeparateParts(left), expr, utils.SeparateParts(right))

	boolVal := "AND"
	if clay {
		boolVal = "OR"
	}
	if q.whereExp != "WHERE (" && q.whereExp != "" && q.whereExp[len(q.whereExp)-1] != '(' {
		statement = fmt.Sprintf(" %s %s", boolVal, statement)
	}

	q.whereExp += statement
	return q
}

func (q *Query[T]) On(col string, expr string, value string) *Query[T] {
	return q.OnStatment(col, expr, value, false)
}

func (q *Query[T]) OrOn(col string, expr string, value string) *Query[T] {
	return q.OnStatment(col, expr, value, true)
}

func (q *Query[T]) OnClause(clause func(*Query[T])) *Query[T] {

	if q.whereExp != "" {
		q.whereExp += " AND "
	}

	q.whereExp += "("
	clause(q)
	q.whereExp += ")"
	return q
}

func (q *Query[T]) OrOnClause(clause func(*Query[T])) *Query[T] {

	if q.whereExp != "" {
		q.whereExp += " OR "
	}

	q.whereExp += "("
	clause(q)
	q.whereExp += ")"
	return q
}

func (q *Query[T]) Limit(n int) *Query[T] {
	q.limitExp = fmt.Sprintf(`LIMIT %d`, n)
	return q
}

func (q *Query[T]) Offset(n int) *Query[T] {
	q.offsetExp = fmt.Sprintf(`OFFSET %d`, n)
	return q
}
