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

	// query.NewQuery(user.model).
	// 	Select([]string{"id", "name"}).
	// 	OrderBy([]string{"id"}, "desc").
	// 	Build().
	// 	Get()

	// query.NewQuery(user.model).
	// 	Select([]string{"id", "name"}).
	// 	Where("id", ">", "1").
	// 	Where("id", "!=", "3").
	// 	Build().
	// 	Get()

	query.NewQuery(user.model).
		From("posts").
		On("posts.id", "=", "users.post_id").
		Where("users.id", "=", "10").
		Where("users.active", "=", "1").
		Where("posts.name", "=", "agartha").
		LeftJoin("categories", func(jc *query.Join) {
			jc.On("categories.id", "=", "posts.category_id")
			// Не используем OR, чтобы не ломать UPDATE
		}).
		// Фильтр на LEFT JOIN таблицу можно оставить, если категория есть
		Where("categories.name", "like", "%A%").
		Update(map[string]any{
			"users.role_id":  1,
			"users.owner_id": 2,
		})

		// println(q1)

		// query.NewQuery(user.model).
		// 	Select([]string{"id", "email"}).
		// 	Where("id", "=", "10").
		// 	Where("active", "=", "1").
		// 	Build().
		// 	SQL()

		// query.NewQuery(user.model).
		// 	Select([]string{"id"}).
		// 	Where("role", "=", "'admin'").
		// 	OrWhere("role", "=", "'moderator'").
		// 	Build().
		// 	SQL()

		// query.NewQuery(user.model).
		// 	Select([]string{"id"}).
		// 	Where("a", "=", "1").
		// 	OrWhere("b", "=", "2").
		// 	Where("c", "=", "3").
		// 	Build().
		// 	SQL()

	query.NewQuery(user.model).
		Select([]string{"id", "name"}).
		Where("status", "=", "1").
		WhereClause(func(q *query.Query) {
			q.
				Where("age", ">", "18").
				OrWhereClause(func(q *query.Query) {
					q.
						Where("role", "=", "'admin'").
						Where("last_login", ">", "'2026-01-01'")
				})
		}).
		Where("active", "=", "1").
		Build().
		SQL()
	// query.NewQuery(user.model).
	// 	Select([]string{"id"}).
	// 	WhereClause(func(q *query.Query) {
	// 		q.
	// 			Where("a", "=", "1").
	// 			OrWhereClause(func(q *query.Query) {
	// 				q.
	// 					Where("b", "=", "2").
	// 					Where("c", "=", "3")
	// 			})
	// 	}).
	// 	Build().
	// 	SQL()

	// query.NewQuery(user.model).
	// 	Select([]string{"users.id", "posts.id"}).
	// 	LeftJoin("posts", "posts.user_id", "users.id").
	// 	Where("users.active", "=", "1").
	// 	WhereClause(func(q *query.Query) {
	// 		q.
	// 			Where("posts.published", "=", "1").
	// 			OrWhere("posts.is_preview", "=", "1")
	// 	}).
	// 	Build().
	// 	SQL()

	// query.NewQuery(user.model).
	// 	Select([]string{"id"}).
	// 	OrderBy([]string{"id"}, "asc").
	// 	Build().
	// 	SQL()

	// query.NewQuery(user.model).
	// 	Create(map[string]string{
	// 		"name":  "'John'",
	// 		"email": "'john@mail.com'",
	// 	})

}
