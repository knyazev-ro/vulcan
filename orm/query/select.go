package query

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan/utils"
)

func (q *Query) Select(cols []string) *Query {
	colsSafe := utils.ColsSafe(cols)
	colsStr := strings.Join(colsSafe, ", ")
	selectStatement := fmt.Sprintf("SELECT %s", colsStr)
	q.selectExp = selectStatement
	return q
}

func (q *Query) From(table string) *Query {
	q.fullStatement += fmt.Sprintf(" FROM %s", table)
	q.fullStatement = strings.Trim(q.fullStatement, " ")
	return q
}

func (q *Query) On(left string, oper string, right string) *Query {
	q.whereExp = fmt.Sprintf("WHERE %s %s %s", utils.SeparateParts(left), oper, utils.SeparateParts(right))
	q.whereExp = strings.Trim(q.whereExp, " ")
	return q
}
