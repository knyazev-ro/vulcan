package tests

import (
	"fmt"

	"github.com/knyazev-ro/vulcan/orm/vulcan"
)

type CommentTest struct {
	_        string `type:"metadata" table:"comments" pk:"id"`
	Id       int64  `type:"column" col:"id"`
	PostId   int64  `type:"column" col:"post_id"`
	Content  string `type:"column" col:"content"`
	Approved int64  `type:"column" col:"approved"`
}
type CategoryTest struct {
	_    string `type:"metadata" table:"categories" pk:"id"`
	Id   int64  `type:"column" col:"id"`
	Name string `type:"column" col:"name"`
}

type PostTest struct {
	_        string        `type:"metadata" table:"posts" pk:"id"`
	Id       int64         `type:"column" col:"id"`
	Name     string        `type:"column" col:"name"`
	UserId   int64         `type:"column" col:"user_id"`
	Category CategoryTest  `type:"relation" table:"categories" reltype:"belongs-to" fk:"category_id" originalkey:"id"`
	Comments []CommentTest `type:"relation" table:"comments" reltype:"has-many" fk:"post_id"`
}

type UserTest struct {
	_        string     `type:"metadata" table:"users" pk:"id"`
	Id       int64      `type:"column" col:"id"`
	Name     string     `type:"column" col:"name"`
	LastName string     `type:"column" col:"last_name"`
	Posts    []PostTest `type:"relation" table:"posts" reltype:"has-many" fk:"user_id"`
}

func ExampleQuery_OrderBy() {
	sql := vulcan.NewQuery[UserTest]().OrderBy("desc", "id").OrderBy("asc", "name", "users.last_name").Build().SQL()
	fmt.Println(sql)
	// Output: SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users ORDER BY "id" DESC, "name", "users"."last_name" ASC;
}

func ExampleQuery_Where() {
	sql := vulcan.NewQuery[UserTest]().
		Where("id", ">", 1).
		Where("id", "!=", 3).
		Build().
		SQL()
	fmt.Println(sql)
	// Output: SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users WHERE "id" > $1 AND "id" != $2;
}

func ExampleQuery_Update() {
	sql := vulcan.NewQuery[UserTest]().
		From("posts").
		On("posts.id", "=", "users.post_id").
		Where("users.id", "=", 10).
		Where("users.active", "=", 1).
		Where("posts.name", "=", "agartha").
		LeftJoin("categories", func(jc *vulcan.Join) {
			jc.On("categories.id", "=", "posts.category_id")
		}).
		Where("categories.name", "like", "%A%").Build().SQL()
	// 	Update(ctx, map[string]any{
	// 		"users.role_id":  1,
	// 		"users.owner_id": 2,
	// 	})
	fmt.Println(sql)
	// Output: SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users FROM "posts" LEFT JOIN "categories" ON "categories"."id" = "posts"."category_id" WHERE "posts"."id" = "users"."post_id" AND "users"."id" = $1 AND "users"."active" = $2 AND "posts"."name" = $3 AND "categories"."name" like $4;
}

func ExampleQuery_OrWhere() {
	sql := vulcan.NewQuery[UserTest]().
		Where("role", "=", "admin").
		OrWhere("role", "=", "moderator").
		Build().
		SQL()
	fmt.Println(sql)
	// Output: SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users WHERE "role" = $1 OR "role" = $2;
}

func ExampleQuery_WhereClause() {
	sql := vulcan.NewQuery[UserTest]().
		InnerJoin("posts", func(jc *vulcan.Join) { jc.On("posts.user_id", "=", "users.id") }).
		LeftJoin("categories", func(jc *vulcan.Join) { jc.On("categories.id", "=", "posts.category_id") }).
		LeftJoin("comments", func(jc *vulcan.Join) { jc.On("comments.post_id", "=", "posts.id") }).
		Where("users.active", "=", 1).
		WhereClause(func(q *vulcan.Query[UserTest]) {
			q.Where("users.status", "=", "premium").
				OrWhereClause(func(q *vulcan.Query[UserTest]) {
					q.Where("users.role", "=", "admin").
						WhereClause(func(q *vulcan.Query[UserTest]) {
							q.Where("users.age", ">", 30).
								OrWhere("users.signup_date", ">", "2025-01-01")
						})
				})
		}).
		Where("posts.published", "=", 1).
		WhereClause(func(q *vulcan.Query[UserTest]) {
			q.Where("categories.name", "like", "%Tech%").
				OrWhere("categories.name", "like", "%Science%")
		}).
		WhereClause(func(q *vulcan.Query[UserTest]) {
			q.Where("comments.approved", "=", 1).
				OrWhere("comments.content", "like", "%important%")
		}).
		Where("posts.views", ">", 1000).
		OrderBy("desc", "users.id", "posts.id").
		Limit(50).
		Offset(10).
		Build().
		SQL()
	fmt.Println(sql)
	// Output: SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users JOIN "posts" ON "posts"."user_id" = "users"."id" LEFT JOIN "categories" ON "categories"."id" = "posts"."category_id" LEFT JOIN "comments" ON "comments"."post_id" = "posts"."id" WHERE "users"."active" = $1 AND ("users"."status" = $2 OR ("users"."role" = $3 AND ("users"."age" > $4 OR "users"."signup_date" > $5))) AND "posts"."published" = $6 AND ("categories"."name" like $7 OR "categories"."name" like $8) AND ("comments"."approved" = $9 OR "comments"."content" like $10) AND "posts"."views" > $11 ORDER BY "users"."id", "posts"."id" DESC LIMIT 50 OFFSET 10;
}
