package vulcan

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan-orm/utils"
)

func (q *Query[T]) GroupBy(cols []string) *Query[T] {
	colsSafe := utils.ColsSafe(cols)
	colsStr := strings.Join(colsSafe, ", ")
	groupByStatement := fmt.Sprintf("GROUP BY %s", colsStr)
	q.groupByExp = groupByStatement
	return q
}
