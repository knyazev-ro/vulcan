package query

import (
	"fmt"
	"strings"
)

func (q *Query) Select(cols []string) *Query {
	colsStr := strings.Join(cols, ", ")
	selectStatement := fmt.Sprintf("SELECT %s", colsStr)
	q.selectExp = selectStatement
	return q
}
