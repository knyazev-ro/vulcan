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
