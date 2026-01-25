package main

import (
	"fmt"
	"time"

	"github.com/knyazev-ro/vulcan/orm/model"
	"github.com/knyazev-ro/vulcan/orm/vulcan"
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

func ExamplesQuery() {

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
		_      string `type:"metadata" table:"posts" pk:"id"`
		Id     int64  `type:"column" col:"id"`
		Name   string `type:"column" col:"name"`
		UserId int64  `type:"column" col:"user_id"`

		// пост принадлежит категории
		Category CategoryTest `type:"relation" table:"categories" reltype:"belongs-to" fk:"category_id" originalkey:"id"`

		// у поста много комментариев
		Comments []CommentTest `type:"relation" table:"comments" reltype:"has-many" fk:"post_id"`
	}

	type UserTest struct {
		_        string     `type:"metadata" table:"users" pk:"id"`
		Id       int64      `type:"column" col:"id"`
		Name     string     `type:"column" col:"name"`
		LastName string     `type:"column" col:"last_name"`
		Posts    []PostTest `type:"relation" table:"posts" reltype:"has-many" fk:"user_id"`
	}

	vulcan.NewQuery[UserTest]().
		OrderBy([]string{"id"}, "desc").
		Build().
		Get()

	vulcan.NewQuery[UserTest]().
		Where("id", ">", 1).
		Where("id", "!=", 3).
		Build().
		Get()

	vulcan.NewQuery[UserTest]().
		From("posts").
		On("posts.id", "=", "users.post_id").
		Where("users.id", "=", 10).
		Where("users.active", "=", 1).
		Where("posts.name", "=", "agartha").
		LeftJoin("categories", func(jc *vulcan.Join) {
			jc.On("categories.id", "=", "posts.category_id")
		}).
		Where("categories.name", "like", "%A%").
		Update(map[string]any{
			"users.role_id":  1,
			"users.owner_id": 2,
		})

	vulcan.NewQuery[UserTest]().
		Where("role", "=", "admin").
		OrWhere("role", "=", "moderator").
		Build().
		SQL()

	vulcan.NewQuery[UserTest]().
		Where("a", "=", 1).
		OrWhere("b", "=", 2).
		Where("c", "=", 3).
		Build().
		SQL()

	vulcan.NewQuery[UserTest]().
		Where("status", "=", 1).
		WhereClause(func(q *vulcan.Query[UserTest]) {
			q.
				Where("age", ">", 18).
				OrWhereClause(func(q *vulcan.Query[UserTest]) {
					q.
						Where("role", "=", "admin").
						Where("last_login", ">", "2026-01-01")
				})
		}).
		Where("active", "=", 1).
		Build().
		SQL()

	vulcan.NewQuery[UserTest]().
		WhereClause(func(q *vulcan.Query[UserTest]) {
			q.
				Where("a", "=", 1).
				OrWhereClause(func(q *vulcan.Query[UserTest]) {
					q.
						Where("b", "=", 2).
						Where("c", "=", 3)
				})
		}).
		Build().
		SQL()

	vulcan.NewQuery[UserTest]().
		OrderBy([]string{"id"}, "asc").
		Build().
		SQL()

	vulcan.NewQuery[UserTest]().
		Create(map[string]any{
			"name":      "John",
			"last_name": "Johanson",
		})

	q := vulcan.NewQuery[UserTest]().
		InnerJoin("posts", func(jc *vulcan.Join) {
			jc.On("posts.user_id", "=", "users.id")
		}).
		LeftJoin("categories", func(jc *vulcan.Join) {
			jc.On("categories.id", "=", "posts.category_id")
		}).
		LeftJoin("comments", func(jc *vulcan.Join) {
			jc.On("comments.post_id", "=", "posts.id")
		}).
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
		OrderBy([]string{"users.id", "posts.id"}, "desc"). // В процессе переработки
		Limit(50).
		Offset(10)

	sql := q.Build().SQL()
	bindings := q.Bindings

	fmt.Println("SQL:", sql)
	fmt.Println("Bindings:", bindings)
}

func ExamplesORM() {

	type TagTest struct {
		_    string `type:"metadata" table:"tags" pk:"id"`
		Id   int64  `type:"column" col:"id"`
		Name string `type:"column" col:"name"`
	}

	type PostTag struct {
		_      string  `type:"metadata" table:"post_tags" pk:"post_id,tag_id" tabletype:"pivot"`
		PostId int64   `type:"column" col:"post_id"`
		TagId  int64   `type:"column" col:"tag_id"`
		Tag    TagTest `type:"relation" table:"tags" reltype:"belongs-to" fk:"tag_id" originalkey:"id"`
	}

	type PostTest struct {
		_        string    `type:"metadata" table:"posts" pk:"id"`
		Id       int64     `type:"column" col:"id"`
		Name     string    `type:"column" col:"name"`
		UserId   int64     `type:"column" col:"user_id"`
		PostTags []PostTag `type:"relation" table:"post_tags" reltype:"has-many" fk:"post_id" originalkey:"id"`
	}

	type ProfileTest struct {
		_      string `type:"metadata" table:"profiles" pk:"id"`
		Id     int64  `type:"column" col:"id"`
		UserId int64  `type:"column" col:"user_id"`
		Bio    string `type:"column" col:"bio"`
		Avatar string `type:"column" col:"avatar"`
	}

	type UserTest struct {
		_        string      `type:"metadata" table:"users" pk:"id"`
		Id       int64       `type:"column" col:"id"`
		Name     string      `type:"column" col:"name"`
		LastName string      `type:"column" col:"last_name"`
		Posts    []PostTest  `type:"relation" table:"posts" reltype:"has-many" fk:"user_id" originalkey:"id"`
		Profile  ProfileTest `type:"relation" table:"profiles" reltype:"has-one" fk:"user_id" originalkey:"id"`
	}

	type DefUserTest struct {
		_        string `type:"metadata" table:"users" pk:"id"`
		Id       int64  `type:"column" col:"id"`
		Name     string `type:"column" col:"name"`
		LastName string `type:"column" col:"last_name"`
	}

	type MainProfileTest struct {
		_      string      `type:"metadata" table:"profiles" pk:"id"`
		Id     int64       `type:"column" col:"id"`
		UserId int64       `type:"column" col:"user_id"`
		Bio    string      `type:"column" col:"bio"`
		Avatar string      `type:"column" col:"avatar"`
		User   DefUserTest `type:"relation" table:"users" reltype:"belongs-to" fk:"user_id" originalkey:"id"`
	}

	// vulcan.NewQuery[UserTest]().Where("name", "like", "Bobby").Update(map[string]any{
	// 	"name":      "Duran",
	// 	"last_name": "Duran",
	// })

	// vulcan.NewQuery[UserTest]().Create(map[string]any{
	// 	"name":      "Garry",
	// 	"last_name": "Debrua",
	// })

	// vulcan.NewQuery[UserTest]().Create(map[string]any{
	// 	"name":      "Bobby",
	// 	"last_name": "Fisher",
	// })

	start := time.Now()
	fmt.Println()
	_ = vulcan.NewQuery[UserTest]().Load()
	// fmt.Println(len(q1))
	end := time.Now()
	fmt.Println(end.Sub(start))

	// model, ok := vulcan.NewQuery[UserTest]().FindById(3)

	// if ok {
	// 	fmt.Println(model)
	// }

	// vulcan.NewQuery[UserTest]().Where("users.name", "like", "%Garry%").Delete()
	// vulcan.NewQuery[UserTest]().DeleteById(1)

	// start = time.Now()
	// fmt.Println()
	// vulcan.NewQuery[UserTest]().
	// 	Build().
	// 	Get()
	// // fmt.Println(q2)
	// end = time.Now()
	// fmt.Println(end.Sub(start))

	// vulcan.NewQuery[UserTest]().Using("posts p", "profiles pr").Where("p.name", "like", "%A%").Delete()
}
