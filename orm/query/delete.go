package query

import (
	"fmt"

	"github.com/knyazev-ro/vulcan/orm/db"
)

func (q *Query[T]) Delete() bool {
	q.fullStatement = fmt.Sprintf(`DELETE FROM %s`, q.Model.TableName)
	q.appendExpressions()
	q.fillBindingsPSQL()
	db := db.DB // предполагаем, что db.DB — это *sql.DB
	res, err := db.Exec(q.fullStatement, q.Bindings...)
	if err != nil {
		panic(err)
	}

	// Опционально: количество удалённых строк
	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deleted %d rows\n", affected)
	return affected > 0
}

func (q *Query[T]) DeleteById(id int64) bool {
	q.Where("id", "=", id)
	q.Delete()
	return true
}
