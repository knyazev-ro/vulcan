package query

import (
	"fmt"
	"strings"
)

func (q *Query) Insert(cols []string, values [][]string) bool {
	valuesStrContainer := []string{}
	for _, val := range values {
		valuesStr := "(" + strings.Join(val, ", ") + ")"
		valuesStrContainer = append(valuesStrContainer, valuesStr)
	}

	valuesStrContainerJoin := strings.Join(valuesStrContainer, ", ")
	colsStr := "(" + strings.Join(cols, ", ") + ")"
	statement := fmt.Sprintf("INSERT INTO %s %s VALUES %s;", q.Model.TableName, colsStr, valuesStrContainerJoin)
	q.fullStatement = statement
	println(statement)
	return true
}

func (q *Query) Create(keyvals map[string]string) *Query {
	values := []string{}
	cols := []string{}
	for k, v := range keyvals {
		cols = append(cols, k)
		values = append(values, v)
	}

	colsStr := "(" + strings.Join(cols, ", ") + ")"
	valuesStr := "(" + strings.Join(values, ", ") + ")"

	statement := fmt.Sprintf("INSERT INTO %s %s VALUES %s;", q.Model.TableName, colsStr, valuesStr)
	q.fullStatement = statement
	fmt.Println(statement)
	// create, then bring back fresh created logic
	// ...
	return q
}
