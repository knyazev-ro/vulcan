package vulcan

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan-orm/utils"
)

func (q *Query[T]) OrderBy(direction string, cols ...string) *Query[T] {
	colsSafe := utils.ColsSafe(cols)
	orderCols := strings.Join(colsSafe, ", ")
	statement := fmt.Sprintf("%s %s", orderCols, strings.ToUpper(direction))
	q.orderExp = append(q.orderExp, statement)
	return q
}
