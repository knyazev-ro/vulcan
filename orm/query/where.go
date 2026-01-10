package query

import (
	"fmt"

	"github.com/knyazev-ro/vulcan/utils"
)

type WhereQuery struct {
}

func (q *Query) WhereStatment(col string, expr string, value string, clay bool) *Query {
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

func (q *Query) Where(col string, expr string, value string) *Query {
	return q.WhereStatment(col, expr, value, false)
}

func (q *Query) OrWhere(col string, expr string, value string) *Query {
	return q.WhereStatment(col, expr, value, true)
}

func (q *Query) WhereClause(clause func(*Query)) *Query {
	if q.whereExp != "" && q.whereExp[len(q.whereExp)-1] != '(' {
		q.whereExp += " AND "
	}

	q.whereExp += "("
	clause(q)
	q.whereExp += ")"
	return q
}

func (q *Query) OrWhereClause(clause func(*Query)) *Query {

	if q.whereExp != "" && q.whereExp[len(q.whereExp)-1] != '(' {
		q.whereExp += " OR "
	}

	q.whereExp += "("
	clause(q)
	q.whereExp += ")"
	return q
}
