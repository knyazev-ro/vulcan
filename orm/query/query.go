package query

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // драйвер для database/sql
	"github.com/knyazev-ro/vulcan/config"
	"github.com/knyazev-ro/vulcan/orm/model"
)

type Query struct {
	Model         model.Model
	Bindings      []any
	selectExp     string
	whereExp      string
	joinExp       string
	createExp     string
	orderExp      string
	fromExp       string
	limitExp      string
	offsetExp     string
	fullStatement string
}

func NewQuery(model model.Model) *Query {
	return &Query{
		Model: model,
	}
}

func (q *Query) Build() *Query {

	if q.fullStatement != "" {
		return q
	}

	selectStr := fmt.Sprintf("SELECT *")

	if q.selectExp != "" {
		selectStr = q.selectExp
	}

	from := fmt.Sprintf("FROM %s", q.Model.TableName)

	q.fullStatement = fmt.Sprintf("%s %s", strings.Trim(selectStr, " "), from)
	q.fullStatement = strings.Trim(q.fullStatement, " ")
	q.appendExpressions()
	q.fullStatement += ";"

	q.fillBindingsPSQL()

	return q
}

func (q *Query) appendExpressions() {

	if q.fromExp != "" {
		q.fullStatement += " " + strings.Trim(q.fromExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.joinExp != "" {
		q.fullStatement += " " + strings.Trim(q.joinExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.whereExp != "" {
		q.fullStatement += " WHERE " + strings.Trim(q.whereExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.orderExp != "" {
		q.fullStatement += " " + strings.Trim(q.orderExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.limitExp != "" {
		q.fullStatement += " " + strings.Trim(q.limitExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}

	if q.offsetExp != "" {
		q.fullStatement += " " + strings.Trim(q.offsetExp, " ")
		q.fullStatement = strings.Trim(q.fullStatement, " ")
	}
}

func (q *Query) SQL() string {
	// println("SQL: ", q.fullStatement)
	return q.fullStatement
}

func (q *Query) RawSQL(v string) *Query {
	q.fullStatement = v
	return q
}

func (q *Query) Get() {
	config := config.GetConfig()
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s", config.Driver, config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("pgx", dsn) // pgx через database/sql
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		panic(err)
	}
	println(q.SQL())
	rows, err := db.Query(q.fullStatement, q.Bindings...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email, role string
		if err := rows.Scan(&id, &name); err != nil {
			panic(err)
		}
		fmt.Println(id, name, email, role)
	}
	return
}

func (q *Query) fillBindingsPSQL() {
	var b strings.Builder
	b.Grow(len(q.fullStatement) + 16)

	idx := 1
	for i := 0; i < len(q.fullStatement); i++ {
		if q.fullStatement[i] == '?' {
			b.WriteString(fmt.Sprintf("$%d", idx))
			idx++
		} else {
			b.WriteByte(q.fullStatement[i])
		}
	}

	q.fullStatement = b.String()
}
