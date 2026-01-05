package query

import "fmt"

func (q *Query) WhereStatment(col string, expr string, value string, clay bool) *Query {
	whereStr := ""
	if q.whereExp == "" {
		whereStr = "WHERE "
	}

	if len(q.whereExp) == 1 && q.whereExp == "(" {
		q.whereExp = "WHERE ("
	}

	statement := fmt.Sprintf(`%s"%s" %s ?`, whereStr, col, expr)
	q.Bindings = append(q.Bindings, value)
	boolVal := "AND"
	if clay {
		boolVal = "OR"
	}
	if q.whereExp != "WHERE (" && q.whereExp != "" && q.whereExp[len(q.whereExp)-1] != '(' {
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

	if q.whereExp != "" {
		q.whereExp += " AND "
	}

	q.whereExp += "("
	clause(q)
	q.whereExp += ")"
	return q
}

func (q *Query) OrWhereClause(clause func(*Query)) *Query {

	if q.whereExp != "" {
		q.whereExp += " OR "
	}

	q.whereExp += "("
	clause(q)
	q.whereExp += ")"
	return q
}
