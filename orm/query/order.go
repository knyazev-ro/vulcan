package query

import (
	"fmt"
	"strings"
)

func (q *Query) OrderBy(cols []string, direction string) *Query {
	orderCols := strings.Join(cols, ", ")
	statement := fmt.Sprintf("ORDER BY %s %s", orderCols, strings.ToUpper(direction))
	q.orderExp = statement
	return q
}
