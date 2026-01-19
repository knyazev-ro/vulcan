package vulcan

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan/orm/db"
)

func (q *Query[T]) Update(values map[string]any) bool {
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
	db := db.DB // предполагаем, что db.DB — это *sql.DB
	res, err := db.Exec(q.fullStatement, q.Bindings...)
	if err != nil {
		panic(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Update affected %d rows\n", affected)

	return true
}
