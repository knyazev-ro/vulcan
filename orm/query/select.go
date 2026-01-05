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
