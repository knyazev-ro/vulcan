package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/knyazev-ro/vulcan/orm/db"
	"github.com/knyazev-ro/vulcan/orm/vulcan"
)

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

func TestVulcanEndurance(t *testing.T) {
	db.Init()
	var wg sync.WaitGroup
	iterations := 500

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := vulcan.NewQuery[ReportData]().CLoad(context.Background())
			if err != nil {
				fmt.Printf("Ошибка в запросе %d: %v\n", id, err)
			}
		}(i)
	}
	wg.Wait()
}
