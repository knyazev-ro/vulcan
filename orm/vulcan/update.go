package vulcan

import (
	"context"
	"fmt"
	"strings"
)

func (q *Query[T]) Update(ctx context.Context, values map[string]any) error {
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
	q.joinExp = ""
	q.fullStatement = fmt.Sprintf(`UPDATE "%s" SET %s`, q.Model.TableName, setsStr) + " " + q.fullStatement
	q.appendExpressions()
	q.fillBindingsPSQL()
	fmt.Println(q.SQL())

	res, err := q.db.ExecContext(ctx, q.fullStatement, q.Bindings...)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("Update affected %d rows\n", affected)

	return nil
}
