package main

import (
	"fmt"
	"time"

	"github.com/knyazev-ro/vulcan/orm/vulcan"
)

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
		Load()

	vulcan.NewQuery[UserTest]().
		Where("id", ">", 1).
		Where("id", "!=", 3).
		Load()

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
	q1 := vulcan.NewQuery[UserTest]().Load()
	end := time.Now()
	fmt.Println(len(q1))
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
	// 	Load()
	// // fmt.Println(q2)
	// end = time.Now()
	// fmt.Println(end.Sub(start))

	// vulcan.NewQuery[UserTest]().Using("posts p", "profiles pr").Where("p.name", "like", "%A%").Delete()
}

func RealExampleORM() {

	type Report struct {
		_  string `type:"metadata" table:"reports" pk:"id"`
		Id int64  `type:"column" col:"id"`
		// Дополнительные поля из fillable
		SystemFileId         int64  `type:"column" col:"system_file_id"`
		DataFilePath         string `type:"column" col:"data_file_path"`
		Filters              string `type:"column" col:"filters"`
		Status               string `type:"column" col:"status"`
		TechnicalTitle       string `type:"column" col:"technical_title"`
		CreatedBy            int64  `type:"column" col:"created_by"`
		ParseErrorMessage    string `type:"column" col:"parse_error_message"`
		FullReportAttachment string `type:"column" col:"full_report_attachment"`
	}

	type ReportIndex struct {
		_             string `type:"metadata" table:"report_indices" pk:"id"`
		Id            int64  `type:"column" col:"id"`
		Index         string `type:"column" col:"index"`
		IndexFullName string `type:"column" col:"indexfullname"`
		IndexCode     string `type:"column" col:"indexcode"`
		IsBlue        bool   `type:"column" col:"is_blue"`
		IsPercentile  bool   `type:"column" col:"is_percentile"`
		IsSalary      bool   `type:"column" col:"is_salary"`
		Order         int64  `type:"column" col:"order"`
		ReportId      int64  `type:"column" col:"report_id"`
	}

	type ReportPeriod struct {
		_        string `type:"metadata" table:"report_periods" pk:"id"`
		Id       int64  `type:"column" col:"id"`
		ReportId int64  `type:"column" col:"report_id"`
		Name     string `type:"column" col:"period"`
	}

	type ReportFuncGroup struct {
		_           string `type:"metadata" table:"report_func_groups" pk:"id"`
		Id          int64  `type:"column" col:"id"`
		ReportId    int64  `type:"column" col:"report_id"`
		Index       string `type:"column" col:"index"`
		Relative    string `type:"column" col:"relative"`
		Text        string `type:"column" col:"text"`
		TypeId      int64  `type:"column" col:"type_id"`
		Description string `type:"column" col:"description"`
		Comments    string `type:"column" col:"comments"`
	}

	type ReportCompGroup struct {
		_             string `type:"metadata" table:"report_comp_groups" pk:"id"`
		Id            int64  `type:"column" col:"id"`
		ReportId      int64  `type:"column" col:"report_id"`
		CompGroup     string `type:"column" col:"compgroup"`
		CompGroupCode string `type:"column" col:"compgroupcode"`
		CompGroupName string `type:"column" col:"compgroupname"`
	}

	type ReportDifficult struct {
		_           string `type:"metadata" table:"report_difficult" pk:"id"`
		Id          int64  `type:"column" col:"id"`
		ReportId    int64  `type:"column" col:"report_id"`
		Code        string `type:"column" col:"code"`
		Name        string `type:"column" col:"name"`
		FuncGroupId int64  `type:"column" col:"func_group_id"`
	}

	type ReportLevel struct {
		_               string `type:"metadata" table:"report_levels" pk:"id"`
		Id              int64  `type:"column" col:"id"`
		ReportId        int64  `type:"column" col:"report_id"`
		Level           string `type:"column" col:"level"`
		LevelCode       string `type:"column" col:"levelcode"`
		LevelName       string `type:"column" col:"levelname"`
		LevelExportName string `type:"column" col:"levelexportname"`
	}

	type ReportCity struct {
		_        string `type:"metadata" table:"report_cities" pk:"id"`
		Id       int64  `type:"column" col:"id"`
		ReportId int64  `type:"column" col:"report_id"`
		City     string `type:"column" col:"city"`
		Order    int64  `type:"column" col:"order"`
	}

	type ReportValueItem struct {
		_           string  `type:"metadata" table:"report_value_items" pk:"id"`
		Id          int64   `type:"column" col:"id"`
		Value       float64 `type:"column" col:"value"`
		IsMoney     bool    `type:"column" col:"is_money"`
		DataTableId int64   `type:"column" col:"data_table_id"`
		LabelId     int64   `type:"column" col:"label_id"`
	}

	type ReportData struct {
		_              string `type:"metadata" table:"report_data" pk:"id"`
		Id             int64  `type:"column" col:"id"`
		ReportId       int64  `type:"column" col:"report_id"`
		Order          string `type:"column" col:"order"`
		CitiesId       int64  `type:"column" col:"cities_id"`
		CompGroupsId   int64  `type:"column" col:"comp_groups_id"`
		DifficultiesId int64  `type:"column" col:"difficulties_id"`
		FuncGroupsId   int64  `type:"column" col:"func_groups_id"`
		IndicatorId    int64  `type:"column" col:"indicator_id"`
		LevelsId       int64  `type:"column" col:"levels_id"`
		PeriodsId      int64  `type:"column" col:"periods_id"`
		QualsId        int64  `type:"column" col:"quals_id"`

		// Relations (belongs-to, однонаправленные)
		Report    Report          `type:"relation" table:"reports" reltype:"belongs-to" fk:"report_id" originalkey:"id"`
		City      ReportCity      `type:"relation" table:"report_cities" reltype:"belongs-to" fk:"cities_id" originalkey:"id"`
		CompGroup ReportCompGroup `type:"relation" table:"report_comp_groups" reltype:"belongs-to" fk:"comp_groups_id" originalkey:"id"`
		Difficult ReportDifficult `type:"relation" table:"report_difficult" reltype:"belongs-to" fk:"difficulties_id" originalkey:"id"`
		FuncGroup ReportFuncGroup `type:"relation" table:"report_func_groups" reltype:"belongs-to" fk:"func_groups_id" originalkey:"id"`
		Level     ReportLevel     `type:"relation" table:"report_levels" reltype:"belongs-to" fk:"levels_id" originalkey:"id"`
		Period    ReportPeriod    `type:"relation" table:"report_periods" reltype:"belongs-to" fk:"periods_id" originalkey:"id"`
		Index     ReportIndex     `type:"relation" table:"report_indices" reltype:"belongs-to" fk:"indicator_id" originalkey:"id"`

		// Дочерние элементы
		ValueItems []ReportValueItem `type:"relation" table:"report_value_items" reltype:"has-many" fk:"data_table_id" originalkey:"id"`
	}

	start := time.Now()
	q2, _ := vulcan.NewQuery[ReportData]().With("City", func(q *vulcan.Query[ReportData]) {
		q.Where("city", "like", "Москва")
	}).FindById(2)
	end := time.Now()
	fmt.Println(q2)
	fmt.Println(end.Sub(start))

	start = time.Now()
	vulcan.NewQuery[ReportData]().Load()
	end = time.Now()
	fmt.Println(end.Sub(start))

}
