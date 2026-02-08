package vulcan

import (
	"context"
	"fmt"

	"github.com/knyazev-ro/vulcan/utils"
)

type WhereQuery struct {
}

func (q *Query[T]) WhereStatment(col string, expr string, value any, clay bool) *Query[T] {
	placeholder := "?"
	statement := fmt.Sprintf(`%s %s %s`, utils.SeparateParts(col), expr, placeholder)
	q.Bindings = append(q.Bindings, value)
	boolVal := "AND"
	if clay {
		boolVal = "OR"
	}
	if q.whereExp != "" && q.whereExp[len(q.whereExp)-1] != '(' {
		statement = fmt.Sprintf(" %s %s", boolVal, statement)
	}

	q.whereExp += statement
	return q
}

func (q *Query[T]) Where(col string, expr string, value any) *Query[T] {
	return q.WhereStatment(col, expr, value, false)
}

func (q *Query[T]) OrWhere(col string, expr string, value any) *Query[T] {
	return q.WhereStatment(col, expr, value, true)
}

func (q *Query[T]) WhereClause(clause func(*Query[T])) *Query[T] {
	if q.whereExp != "" && q.whereExp[len(q.whereExp)-1] != '(' {
		q.whereExp += " AND "
	}

	q.whereExp += "("
	clause(q)
	q.whereExp += ")"
	return q
}

func (q *Query[T]) OrWhereClause(clause func(*Query[T])) *Query[T] {

	if q.whereExp != "" && q.whereExp[len(q.whereExp)-1] != '(' {
		q.whereExp += " OR "
	}

	q.whereExp += "("
	clause(q)
	q.whereExp += ")"
	return q
}

func (q *Query[T]) WhereAny(col string, values []any) *Query[T] {

	statement := fmt.Sprintf(`%s = ANY(?)`, col)
	q.Bindings = append(q.Bindings, values)
	boolVal := "AND"
	if q.whereExp != "" && q.whereExp[len(q.whereExp)-1] != '(' {
		statement = fmt.Sprintf(" %s %s", boolVal, statement)
	}

	q.whereExp += statement
	return q
}

// return Model and bool - true - exists, false - not
func (q *Query[T]) FindById(ctx context.Context, ids ...int64) (T, bool, error) {
	pks := q.Model.Pks
	var empty T

	if len(pks) != len(ids) {
		return empty, false, &FindByIdError{message: "len of ids arguments and len of pk mismatch. len of pks should be == to len of ids arguments!"}
	}

	for idx, pk := range pks {
		Id := fmt.Sprintf("%s.%s", q.Model.TableName, pk)
		q.Where(Id, "=", ids[idx])
	}

	Ts, err := q.Load(ctx)
	if err != nil {
		return empty, false, err
	}
	if len(Ts) >= 1 {
		return Ts[0], true, nil
	}
	return empty, false, nil
}
