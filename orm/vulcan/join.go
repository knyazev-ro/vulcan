package vulcan

import (
	"fmt"

	"github.com/knyazev-ro/vulcan-orm/utils"
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

func (q *Query[T]) typeJoin(typeJoinName string, table string, clause func(jc *Join)) *Query[T] {
	joinStr := fmt.Sprintf(`%s "%s" ON`, typeJoinName, table)
	join := &Join{body: joinStr}
	clause(join)
	q.joinExp += fmt.Sprintf(` %s %s`, join.body, join.ons)

	return q
}

func (q *Query[T]) LeftJoin(table string, clause func(jc *Join)) *Query[T] {
	q.typeJoin("LEFT JOIN", table, clause)
	return q
}

func (q *Query[T]) RightJoin(table string, clause func(jc *Join)) *Query[T] {
	q.typeJoin("RIGHT JOIN", table, clause)
	return q
}

func (q *Query[T]) InnerJoin(table string, clause func(jc *Join)) *Query[T] {
	q.typeJoin("JOIN", table, clause)
	return q
}
