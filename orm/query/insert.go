package query

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan/utils"
)

func (q *Query) Insert(cols []string, values [][]any) bool {
	valuesStrContainer := []string{}
	for _, val := range values {
		valQ := []string{}
		for range val {
			valQ = append(valQ, "?")
		}
		valuesStr := "(" + strings.Join(valQ, ", ") + ")"
		q.Bindings = append(q.Bindings, val...)
		valuesStrContainer = append(valuesStrContainer, valuesStr)
	}
	valuesStrContainerJoin := strings.Join(valuesStrContainer, ", ")
	colsSafe := utils.ColsSafe(cols)
	colsStr := "(" + strings.Join(colsSafe, ", ") + ")"
	statement := fmt.Sprintf(`INSERT INTO %s %s VALUES %s;`, q.Model.TableName, colsStr, valuesStrContainerJoin)
	q.fullStatement = statement
	q.fillBindingsPSQL()
	println(q.fullStatement)
	return true
}

func ValuesFilledWithQuestions(values []string, q *Query) []string {
	questions := []string{}
	for _, v := range values {
		q.Bindings = append(q.Bindings, v)
		// index := fmt.Sprintf("$%d", i+1)
		index := "?"
		questions = append(questions, index)
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
	q.fillBindingsPSQL()
	fmt.Println(q.fullStatement, q.Bindings)
	// create, then bring back fresh created logic
	// ...
	return q
}
