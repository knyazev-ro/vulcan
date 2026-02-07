package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/knyazev-ro/vulcan/orm/db"
	"github.com/knyazev-ro/vulcan/orm/vulcan"
)

func ExamplesORM() {
	db.Init()
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
	ctx := context.Background()

	err := vulcan.NewQuery[UserTest]().
		From("posts", "tags").
		On("posts.user_id", "=", "users.id"). // Связь User -> Posts
		Where("users.id", "=", 3).
		Where("posts.name", "=", "Zachary Terminals: Budget Edition").
		LeftJoin("tags", func(jc *vulcan.Join) {
			jc.On("tags.id", "=", "posts.id") // Пример джоина для фильтрации
		}).
		Where("tags.name", "like", "Hardware").
		Update(ctx, map[string]any{
			"name":      "Deadman",
			"last_name": "Surree",
		})
	if err != nil {
		fmt.Println("Difficult update error: ", err.Error())
	}

	zack, ok, err := vulcan.NewQuery[UserTest]().FindById(ctx, 3)
	if err != nil {
		fmt.Println("Error during find: ", err.Error())
	}
	if ok {
		fmt.Println(zack.Name, zack.LastName)
	}

	// Update
	err = vulcan.NewQuery[UserTest]().
		Where("name", "like", "Bobby").
		Update(ctx, map[string]any{
			"name":      "Duran",
			"last_name": "Duran",
		})
	if err != nil {
		fmt.Println("Update error:", err)
	} else {
		fmt.Println("Update succeeded")
	}

	// Create новые записи
	err = vulcan.NewQuery[UserTest]().Create(ctx, map[string]any{
		"name":      "Garry",
		"last_name": "Debrua",
	})
	if err != nil {
		fmt.Println("Create error:", err)
	} else {
		fmt.Println("Created user Garry")
	}

	err = vulcan.NewQuery[UserTest]().Create(ctx, map[string]any{
		"name":      "Bobby",
		"last_name": "Fisher",
	})
	if err != nil {
		fmt.Println("Create error:", err)
	} else {
		fmt.Println("Created user Bobby")
	}

	// Load всех пользователей
	start := time.Now()
	users, err := vulcan.NewQuery[UserTest]().CLoad(ctx)
	end := time.Now()
	if err != nil {
		fmt.Println("Load error:", err)
	} else {
		fmt.Println("Duration:", end.Sub(start))
		fmt.Println("Loaded users:", len(users))
	}

	// FindById
	user, ok, err := vulcan.NewQuery[UserTest]().FindById(ctx, 3)
	if err != nil {
		fmt.Println("FindById error:", err)
	} else if ok {
		fmt.Println("User found:", user)
	} else {
		fmt.Println("User not found")
	}

	// Delete по условию
	// _, err = vulcan.NewQuery[UserTest]().
	// 	Where("name", "like", "%Garry%").
	// 	Delete(ctx)
	// if err != nil {
	// 	fmt.Println("Delete error:", err)
	// } else {
	// 	fmt.Println("Deleted users with name like Garry")
	// }

	// // DeleteById
	// _, err = vulcan.NewQuery[UserTest]().DeleteById(ctx, 1)
	// if err != nil {
	// 	fmt.Println("DeleteById error:", err)
	// } else {
	// 	fmt.Println("Deleted user with ID 1")
	// }

	// Delete с Using (пример для множественных таблиц)
	// _, err = vulcan.NewQuery[UserTest]().
	// 	Using("posts p", "profiles pr").
	// 	Where("p.name", "like", "%A%").
	// 	Delete(ctx)
	// if err != nil {
	// 	fmt.Println("Delete with Using error:", err)
	// } else {
	// 	fmt.Println("Deleted users with posts.name like %A%")
	// }
}

func RealExampleORM() {
	db.Init()
	type ReportWithAggCount struct {
		_      string `type:"metadata" table:"reports" pk:"id"`
		Status string `type:"column" col:"status"`
		Count  int64  `type:"column" col:"*" agg:"count"`
	}

	type Report struct {
		_  string `type:"metadata" table:"reports" pk:"id"`
		Id int64  `type:"column" col:"id"`
		// Дополнительные поля из fillable
		SystemFileId         int64   `type:"column" col:"system_file_id"`
		DataFilePath         string  `type:"column" col:"data_file_path"`
		Filters              string  `type:"column" col:"filters"`
		Status               string  `type:"column" col:"status"`
		TechnicalTitle       string  `type:"column" col:"technical_title"`
		CreatedBy            int64   `type:"column" col:"created_by"`
		ParseErrorMessage    *string `type:"column" col:"parse_error_message"`
		FullReportAttachment *string `type:"column" col:"full_report_attachment"`
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

	ctx := context.Background()
	// runtime.GOMAXPROCS(1)
	// FindById с фильтром
	report, ok, err := vulcan.NewQuery[ReportData]().With("City", func(q *vulcan.Query[ReportData]) {
		q.Where("city", "like", "Москва")
	}).FindById(ctx, 2)
	if err != nil {
		fmt.Println("FindById error:", err)
	} else if ok {
		fmt.Println("Report found:", report)
	} else {
		fmt.Println("Report not found")
	}

	// Load no gorutines in transaction
	fmt.Println("SYNC LOAD! USE 'LOAD()'")
	start := time.Now()
	var reports []ReportData
	err = vulcan.Transaction(ctx, func(tx *sql.Tx) error {
		query, err := vulcan.NewQuery[ReportData]().UseConn(tx)
		if err != nil {
			return err
		}
		reports, err = query.Load(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	end := time.Now()
	if err != nil {
		fmt.Println("Load error:", err)
	}
	fmt.Println("Duration:", end.Sub(start))
	fmt.Println("Loaded reports:", len(reports))

	// Load всех записей
	fmt.Println("CONCURRENT LOAD! USE 'CLOAD()'")

	start = time.Now()
	reports, err = vulcan.NewQuery[ReportData]().CLoad(ctx)
	end = time.Now()
	if err != nil {
		fmt.Println("Load error:", err)
	}
	fmt.Println("Duration:", end.Sub(start))
	fmt.Println("Loaded reports:", len(reports))

	// fmt.Println(vulcan.NewQuery[ReportWithAggCount]().CLoad())

}
