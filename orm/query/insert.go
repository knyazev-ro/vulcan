package query

import (
	"fmt"
	"strings"

	"github.com/knyazev-ro/vulcan/orm/db"
	"github.com/knyazev-ro/vulcan/utils"
)

func (q *Query[T]) Insert(cols []string, values [][]any) bool {
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

	db := db.DB // предполагаем, что db.DB — это *sql.DB
	res, err := db.Exec(q.fullStatement, q.Bindings...)
	if err != nil {
		panic(err)
	}

	// Опционально: количество вставленных строк
	affected, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Inserted %d rows\n", affected)
	return affected > 0
}

func (q *Query[T]) valuesFilledWithQuestions(values []any) []string {
	questions := []string{}
	for _, v := range values {
		q.Bindings = append(q.Bindings, v)
		// index := fmt.Sprintf("$%d", i+1)
		index := "?"
		questions = append(questions, index)
	}
	return questions
}

func (q *Query[T]) Create(keyvals map[string]any) *Query[T] {
	values := []any{}
	cols := []string{}
	for k, v := range keyvals {
		cols = append(cols, k)
		values = append(values, v)
	}

	colsSafe := utils.ColsSafe(cols)
	colsStr := "(" + strings.Join(colsSafe, ", ") + ")"

	valuesStr := "(" + strings.Join(q.valuesFilledWithQuestions(values), ", ") + ")"

	statement := fmt.Sprintf(`INSERT INTO %s %s VALUES %s;`, q.Model.TableName, colsStr, valuesStr)
	q.fullStatement = statement
	q.fillBindingsPSQL()

	db := db.DB // предполагаем, что db.DB — это *sql.DB
	res, err := db.Exec(q.fullStatement, q.Bindings...)
	if err != nil {
		panic(err)
	}

	lastID, err := res.LastInsertId()
	if err == nil {
		fmt.Printf("Inserted row ID: %d\n", lastID)
	}

	return q
}
