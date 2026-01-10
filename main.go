package main

import (
	"fmt"

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

	query.NewQuery(user.model).
		Select([]string{"id", "name"}).
		OrderBy([]string{"id"}, "desc").
		Build().
		Get()

	query.NewQuery(user.model).
		Select([]string{"id", "name"}).
		Where("id", ">", "1").
		Where("id", "!=", "3").
		Build().
		Get()

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

	query.NewQuery(user.model).
		Select([]string{"id", "email"}).
		Where("id", "=", "10").
		Where("active", "=", "1").
		Build().
		SQL()

	query.NewQuery(user.model).
		Select([]string{"id"}).
		Where("role", "=", "'admin'").
		OrWhere("role", "=", "'moderator'").
		Build().
		SQL()

	query.NewQuery(user.model).
		Select([]string{"id"}).
		Where("a", "=", "1").
		OrWhere("b", "=", "2").
		Where("c", "=", "3").
		Build().
		SQL()

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

	query.NewQuery(user.model).
		Select([]string{"id"}).
		WhereClause(func(q *query.Query) {
			q.
				Where("a", "=", "1").
				OrWhereClause(func(q *query.Query) {
					q.
						Where("b", "=", "2").
						Where("c", "=", "3")
				})
		}).
		Build().
		SQL()

	query.NewQuery(user.model).
		Select([]string{"id"}).
		OrderBy([]string{"id"}, "asc").
		Build().
		SQL()

	query.NewQuery(user.model).
		Create(map[string]string{
			"name":  "'John'",
			"email": "'john@mail.com'",
		})

	// q := query.NewQuery(user.model).
	// 	Select([]string{"users.id", "users.name", "posts.title", "categories.name"}).
	// 	InnerJoin("posts", func(jc *query.Join) {
	// 		jc.On("posts.user_id", "=", "users.id")
	// 	}).
	// 	LeftJoin("categories", func(jc *query.Join) {
	// 		jc.On("categories.id", "=", "posts.category_id")
	// 	}).
	// 	Where("users.active", "=", "1").
	// 	WhereClause(func(q *query.Query) {
	// 		q.
	// 			Where("users.status", "=", "'premium'").
	// 			OrWhereClause(func(q *query.Query) {
	// 				q.
	// 					Where("users.role", "=", "'admin'").
	// 					WhereClause(func(q *query.Query) {
	// 						q.
	// 							Where("users.age", ">", "30").
	// 							OrWhere("users.signup_date", ">", "'2025-01-01'")
	// 					})
	// 			})
	// 	}).
	// 	Where("posts.published", "=", "1").
	// 	WhereClause(func(q *query.Query) {
	// 		q.
	// 			Where("categories.name", "like", "'%Tech%'").
	// 			OrWhere("categories.name", "like", "'%Science%'")
	// 	}).
	// 	OrderBy([]string{"users.id"}, "desc").
	// 	Limit(50)

	// sql := q.Build().SQL()
	// bindings := q.Bindings

	// fmt.Println("SQL:", sql)
	// fmt.Println("Bindings:", bindings)

	q := query.NewQuery(user.model). // users
						Select([]string{"users.id", "users.name", "posts.title", "categories.name", "comments.content"}).
						InnerJoin("posts", func(jc *query.Join) {
			jc.On("posts.user_id", "=", "users.id")
		}).
		LeftJoin("categories", func(jc *query.Join) {
			jc.On("categories.id", "=", "posts.category_id")
		}).
		LeftJoin("comments", func(jc *query.Join) {
			jc.On("comments.post_id", "=", "posts.id")
		}).
		Where("users.active", "=", "1").
		WhereClause(func(q *query.Query) {
			q.Where("users.status", "=", "'premium'").
				OrWhereClause(func(q *query.Query) {
					q.Where("users.role", "=", "'admin'").
						WhereClause(func(q *query.Query) {
							q.Where("users.age", ">", "30").
								OrWhere("users.signup_date", ">", "'2025-01-01'")
						})
				})
		}).
		Where("posts.published", "=", "1").
		WhereClause(func(q *query.Query) {
			q.Where("categories.name", "like", "'%Tech%'").
				OrWhere("categories.name", "like", "'%Science%'")
		}).
		WhereClause(func(q *query.Query) {
			q.Where("comments.approved", "=", "1").
				OrWhere("comments.content", "like", "'%important%'")
		}).
		Where("posts.views", ">", "1000").
		OrderBy([]string{"users.id", "posts.id"}, "desc").
		Limit(50).
		Offset(10)

	sql := q.Build().SQL()
	bindings := q.Bindings

	fmt.Println("SQL:", sql)
	fmt.Println("Bindings:", bindings)

}
