package query

import "fmt"

func (q *Query) LeftJoin(table string, left string, right string) *Query {
	joinStr := fmt.Sprintf(` LEFT JOIN "%s" ON "%s" = "%s"`, table, left, right)
	q.joinExp += joinStr
	return q
}

func (q *Query) RightJoin(table string, left string, right string) *Query {
	joinStr := fmt.Sprintf(` RIGHT JOIN "%s" ON "%s" = "%s"`, table, left, right)
	q.joinExp += joinStr
	return q
}

func (q *Query) InnerJoin(table string, left string, right string) *Query {
	joinStr := fmt.Sprintf(` JOIN "%s" ON "%s" = "%s"`, table, left, right)
	q.joinExp += joinStr
	return q
}
