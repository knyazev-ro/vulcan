package vulcan

import (
	"context"
	"database/sql"

	"github.com/knyazev-ro/vulcan/orm/db"
)

type DBConnection interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func (q *Query[T]) UseConn(db DBConnection) (*Query[T], error) {
	q.db = db
	return q, nil
}

func (q *Query[T]) Transaction(ctx context.Context, closure func(tx *sql.Tx) error) error {

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = closure(tx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
