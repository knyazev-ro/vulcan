package query

import (
	"fmt"
	"strings"
)

func (q *Query) Update(values map[string]any) bool {
	sets := []string{}
	binds := []any{}
	for i, v := range values {
		parts := strings.Split(i, ".")
		secPart := ""
		if len(parts) == 2 {
			secPart = fmt.Sprintf(`."%s"`, parts[1])
		}
		currSet := fmt.Sprintf(`"%s"%s = ?`, parts[0], secPart)
		sets = append(sets, currSet)
		binds = append(binds, v)
	}

	q.Bindings = append(binds, q.Bindings...)

	setsStr := strings.Join(sets, ", ")
	q.fullStatement = fmt.Sprintf(`UPDATE "%s" SET %s`, q.Model.TableName, setsStr) + " " + q.fullStatement
	q.appendExpressions()
	q.fillBindingsPSQL()
	println(q.SQL())
	// logic
	return true
}
