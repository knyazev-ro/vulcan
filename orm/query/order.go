package query

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan/utils"
)

func (q *Query) OrderBy(cols []string, direction string) *Query {
	colsSafe := utils.ColsSafe(cols)
	orderCols := strings.Join(colsSafe, ", ")
	statement := fmt.Sprintf("ORDER BY %s %s", orderCols, strings.ToUpper(direction))
	q.orderExp = statement
	return q
}
