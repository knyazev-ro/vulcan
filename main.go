package main

import (
	"github.com/knyazev-ro/vulcan/orm/model"
	"github.com/knyazev-ro/vulcan/orm/query"
)

type User struct {
	model model.Model
	Id    int `db:"Id"`
}

func NewUser() *User {
	return &User{
		model: model.Model{
			TableName: "users",
		},
	}
}

func main() {
	user := NewUser()
	q := query.NewQuery(user.model).
		Select([]string{"post_id", "id", "name", "last_name"}).
		Where("id", "=", "25").
		Where("id", "<", "30").
		WhereClause(func(q *query.Query) {
			q.
				Where("id", "<", "15").
				OrWhere("id", ">", "85").
				Where("id", ">", "66").
				OrWhere("id", ">", "85").
				OrWhereClause(func(quer *query.Query) {
					quer.Where("id", ">", "6567")
				})
		}).
		LeftJoin("posts", "posts.id", "users.post_id").
		OrderBy([]string{"id"}, "asc").
		Build().
		SQL()
	println(q)

	query.NewQuery(user.model).Insert([]string{"col1", "col2", "col3"}, [][]string{
		{"1", "2", "3"},
		{"4", "5", "6"},
		{"7", "8", "9"},
		{"10", "11", "12"},
	})

	query.NewQuery(user.model).Create(map[string]string{
		"col1": "1",
		"col2": "2",
		"col3": "3",
	})
}
