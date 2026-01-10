package query

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan/utils"
)

func (q *Query) Select(cols []string) *Query {
	colsSafe := utils.ColsSafe(cols)
	colsStr := strings.Join(colsSafe, ", ")
	selectStatement := fmt.Sprintf("SELECT %s", colsStr)
	q.selectExp = selectStatement
	return q
}

func (q *Query) From(table string) *Query {
	q.fromExp = fmt.Sprintf(`FROM "%s"`, table)
	return q
}

func (q *Query) OnStatment(left string, expr string, right string, clay bool) *Query {
	whereStr := ""
	if q.whereExp == "" {
		whereStr = "WHERE "
	}

	if len(q.whereExp) == 1 && q.whereExp == "(" {
		q.whereExp = "WHERE ("
	}

	statement := fmt.Sprintf(`%s %s %s %s`, whereStr, utils.SeparateParts(left), expr, utils.SeparateParts(right))

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

func (q *Query) On(col string, expr string, value string) *Query {
	return q.OnStatment(col, expr, value, false)
}

func (q *Query) OrOn(col string, expr string, value string) *Query {
	return q.OnStatment(col, expr, value, true)
}

func (q *Query) OnClause(clause func(*Query)) *Query {

	if q.whereExp != "" {
		q.whereExp += " AND "
	}

	q.whereExp += "("
	clause(q)
	q.whereExp += ")"
	return q
}

func (q *Query) OrOnClause(clause func(*Query)) *Query {

	if q.whereExp != "" {
		q.whereExp += " OR "
	}

	q.whereExp += "("
	clause(q)
	q.whereExp += ")"
	return q
}
