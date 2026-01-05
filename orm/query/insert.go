package query

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan/utils"
)

func (q *Query) Insert(cols []string, values [][]string) bool {
	valuesStrContainer := []string{}
	for _, val := range values {
		valuesStr := "(" + strings.Join(val, ", ") + ")"
		valuesStrContainer = append(valuesStrContainer, valuesStr)
	}
	valuesStrContainerJoin := strings.Join(valuesStrContainer, ", ")
	colsSafe := utils.ColsSafe(cols)
	colsStr := "(" + strings.Join(colsSafe, ", ") + ")"
	statement := fmt.Sprintf(`INSERT INTO %s %s VALUES %s;`, q.Model.TableName, colsStr, valuesStrContainerJoin)
	q.fullStatement = statement
	println(statement)
	return true
}

func ValuesFilledWithQuestions(values []string, q *Query) []string {
	questions := []string{}
	for i, v := range values {
		q.Bindings = append(q.Bindings, v)
		questions = append(questions, fmt.Sprintf("$%d", i+1))
	}
	return questions
}

func (q *Query) Create(keyvals map[string]string) *Query {
	values := []string{}
	cols := []string{}
	for k, v := range keyvals {
		cols = append(cols, k)
		values = append(values, v)
	}

	colsSafe := utils.ColsSafe(cols)
	colsStr := "(" + strings.Join(colsSafe, ", ") + ")"

	valuesStr := "(" + strings.Join(ValuesFilledWithQuestions(values, q), ", ") + ")"

	statement := fmt.Sprintf(`INSERT INTO %s %s VALUES %s;`, q.Model.TableName, colsStr, valuesStr)
	q.fullStatement = statement
	fmt.Println(statement, q.Bindings)
	// create, then bring back fresh created logic
	// ...
	return q
}
