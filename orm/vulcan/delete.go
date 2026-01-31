package vulcan

import (
	"context"
	"fmt"

	"github.com/knyazev-ro/vulcan/orm/db"
)

func (q *Query[T]) Delete(ctx context.Context) (int64, error) {
	q.joinExp = ""
	q.fullStatement = fmt.Sprintf(`DELETE FROM %s`, q.Model.TableName)
	q.appendExpressions()
	q.fillBindingsPSQL()
	db := db.DB // предполагаем, что db.DB — это *sql.DB
	res, err := db.ExecContext(ctx, q.fullStatement, q.Bindings...)
	if err != nil {
		return 0, err
	}

	// Опционально: количество удалённых строк
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	fmt.Printf("Deleted %d rows\n", affected)
	return affected, nil
}

func (q *Query[T]) DeleteById(ctx context.Context, id int64) (bool, error) {
	q.Where("id", "=", id)
	aff, err := q.Delete(ctx)
	if err != nil {
		return false, err
	}
	return aff > 0, nil
}
