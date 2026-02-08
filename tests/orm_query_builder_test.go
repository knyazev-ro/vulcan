package tests

import (
	"fmt"
	"testing"

	"github.com/knyazev-ro/vulcan/orm/vulcan"
)

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

func ExampleQuery_OrderBy() {
	sql := vulcan.NewQuery[UserTest]().
		OrderBy("desc", "id").
		OrderBy("asc", "name", "users.last_name").
		Build().SQL()

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
	// Демонстрация сложного Update с Join и подгрузкой из другой таблицы (posts)
	// db.Init()
	sql := vulcan.NewQuery[UserTest]().
		From("posts", "tags").
		On("posts.user_id", "=", "users.id"). // Связь User -> Posts
		Where("users.id", "=", 3).
		Where("posts.name", "=", "Zachary Terminals: Budget Edition").
		LeftJoin("tags", func(jc *vulcan.Join) {
			jc.On("tags.id", "=", "posts.id") // Пример джоина для фильтрации
		}).
		Where("tags.name", "like", "Hardware")
	// Update(ctx, map[string]any{
	// 	"name":      "Deadman",
	// 	"last_name": "Surree",
	// })

	fmt.Println(sql)
}

func ExampleQuery_OrWhere() {
	sql := vulcan.NewQuery[UserTest]().
		Where("name", "=", "Luka").
		OrWhere("name", "=", "Zachary").
		Build().
		SQL()

	fmt.Println(sql)
	// Output: SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users WHERE "name" = $1 OR "name" = $2;
}

func ExampleQuery_WhereClause() {
	// Самый сложный тест: вложенные условия и множественные джоины
	sql := vulcan.NewQuery[UserTest]().
		InnerJoin("posts", func(jc *vulcan.Join) {
			jc.On("posts.user_id", "=", "users.id")
		}).
		LeftJoin("profiles", func(jc *vulcan.Join) {
			jc.On("profiles.user_id", "=", "users.id")
		}).
		Where("users.id", ">", 0).
		WhereClause(func(q *vulcan.Query[UserTest]) {
			q.Where("users.name", "=", "Luka").
				OrWhereClause(func(q *vulcan.Query[UserTest]) {
					q.Where("users.last_name", "=", "Original").
						WhereClause(func(q *vulcan.Query[UserTest]) {
							q.Where("profiles.bio", "like", "%soul%").
								OrWhere("profiles.avatar", "=", "luka_clone.png")
						})
				})
		}).
		Where("posts.name", "like", "%Memory%").
		OrderBy("desc", "users.id", "posts.id").
		Limit(50).
		Offset(10).
		Build().
		SQL()

	fmt.Println(sql)
	// Output: SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users JOIN "posts" ON "posts"."user_id" = "users"."id" LEFT JOIN "profiles" ON "profiles"."user_id" = "users"."id" WHERE "users"."id" > $1 AND ("users"."name" = $2 OR ("users"."last_name" = $3 AND ("profiles"."bio" like $4 OR "profiles"."avatar" = $5))) AND "posts"."name" like $6 ORDER BY "users"."id", "posts"."id" DESC LIMIT 50 OFFSET 10;
}

func TestQuery_Pagination(t *testing.T) {
	t.Run("check sql query for 1 page", func(t *testing.T) {
		sql := vulcan.NewQuery[UserTest]().Paginate("id", 1, 10).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users ORDER BY "id" ASC LIMIT 10 OFFSET 0;` {
			t.Errorf("pagination test failed page 1")
		}
	})

	t.Run("check sql query for 2 page", func(t *testing.T) {
		sql := vulcan.NewQuery[UserTest]().Paginate("id", 2, 10).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users ORDER BY "id" ASC LIMIT 10 OFFSET 10;` {
			t.Errorf("pagination test failed page 2")
		}
	})

	t.Run("check sql query for 3 page", func(t *testing.T) {
		sql := vulcan.NewQuery[UserTest]().Paginate("id", 3, 10).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users ORDER BY "id" ASC LIMIT 10 OFFSET 20;` {
			t.Errorf("pagination test failed page 3")
		}
	})

	t.Run("check sql query for page 1 perpage 5", func(t *testing.T) {
		sql := vulcan.NewQuery[UserTest]().Paginate("id", 1, 5).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users ORDER BY "id" ASC LIMIT 5 OFFSET 0;` {
			t.Errorf("pagination test failed page 3 perpage 5")
		}
	})

	t.Run("check sql query for page 2 perpage 5", func(t *testing.T) {
		sql := vulcan.NewQuery[UserTest]().Paginate("id", 2, 5).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users ORDER BY "id" ASC LIMIT 5 OFFSET 5;` {
			t.Errorf("pagination test failed page 3 perpage 5")
		}
	})
	t.Run("check sql query for page 3 perpage 5", func(t *testing.T) {
		sql := vulcan.NewQuery[UserTest]().Paginate("id", 3, 5).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users ORDER BY "id" ASC LIMIT 5 OFFSET 10;` {
			t.Errorf("pagination test failed page 3 perpage 5")
		}
	})

	t.Run("check sql query for page 128 perpage 100", func(t *testing.T) {
		sql := vulcan.NewQuery[UserTest]().Paginate("id", 128, 100).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users ORDER BY "id" ASC LIMIT 100 OFFSET 12700;` {
			t.Errorf("pagination test failed page 128 perpage 100")
		}
	})

	t.Run("check sql query for page 128 perpage 100 column created_at", func(t *testing.T) {
		sql := vulcan.NewQuery[UserTest]().Paginate("created_at", 128, 100).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users ORDER BY "created_at" ASC LIMIT 100 OFFSET 12700;` {
			t.Errorf("pagination test failed page 128 perpage 100")
		}
	})
}

func TestQuery_CursorPagination(t *testing.T) {
	t.Run("check cursor pagination for 1 page", func(t *testing.T) {
		sql := vulcan.NewQuery[UserTest]().CursorPaginate("id", nil, 10).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users ORDER BY "id" ASC LIMIT 10;` {
			t.Errorf("cursor pagination failed")
		}
	})
	t.Run("check cursor pagination for 2 page", func(t *testing.T) {
		afterLastOne := 14
		sql := vulcan.NewQuery[UserTest]().CursorPaginate("id", afterLastOne, 10).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users WHERE "id" > $1 ORDER BY "id" ASC LIMIT 10;` {
			t.Errorf("cursor pagination failed for page 2 cursor by id 14")
		}
	})

	t.Run("check cursor pagination for 2 page change limit", func(t *testing.T) {
		afterLastOne := 14
		sql := vulcan.NewQuery[UserTest]().CursorPaginate("id", afterLastOne, 25).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users WHERE "id" > $1 ORDER BY "id" ASC LIMIT 25;` {
			t.Errorf("cursor pagination failed for page 2 cursor by id 14 perpage 25")
		}
	})

	t.Run("check cursor pagination limit 25 use different column", func(t *testing.T) {
		sql := vulcan.NewQuery[UserTest]().CursorPaginate("created_at", "2026-02-08 10:30:00", 25).Build().SQL()
		fmt.Println(sql)
		if sql != `SELECT "users"."id" AS users_id, "users"."name" AS users_name, "users"."last_name" AS users_last_name FROM users WHERE "created_at" > $1 ORDER BY "created_at" ASC LIMIT 25;` {
			t.Errorf("cursor pagination failed for perpage 25 and use created_at")
		}
	})
}
