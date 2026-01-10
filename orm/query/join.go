package query

import (
	"fmt"

	"github.com/knyazev-ro/vulcan/utils"
)

type Join struct {
	body string
	ons  string
}

func (j *Join) On(left string, op string, right string) *Join {
	on := fmt.Sprintf(`%s %s %s`, utils.SeparateParts(left), op, utils.SeparateParts(right))
	if j.ons == "" {
		j.ons = on
	} else {
		j.ons += fmt.Sprintf(` AND %s`, on)
	}
	return j
}

func (j *Join) OrOn(left string, op string, right string) *Join {
	on := fmt.Sprintf(`%s %s %s`, utils.SeparateParts(left), op, utils.SeparateParts(right))
	if j.ons == "" {
		j.ons = on
	} else {
		j.ons += fmt.Sprintf(` OR %s`, on)
	}
	return j
}

func (q *Query) typeJoin(typeJoinName string, table string, clause func(jc *Join)) *Query {
	joinStr := fmt.Sprintf(`%s "%s" ON`, typeJoinName, table)
	join := &Join{body: joinStr}
	clause(join)
	q.joinExp += fmt.Sprintf(` %s %s`, join.body, join.ons)

	return q
}

func (q *Query) LeftJoin(table string, clause func(jc *Join)) *Query {
	q.typeJoin("LEFT JOIN", table, clause)
	return q
}

func (q *Query) RightJoin(table string, clause func(jc *Join)) *Query {
	q.typeJoin("RIGHT JOIN", table, clause)
	return q
}

func (q *Query) InnerJoin(table string, clause func(jc *Join)) *Query {
	q.typeJoin("JOIN", table, clause)
	return q
}
