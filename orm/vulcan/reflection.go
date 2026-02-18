package vulcan

import (
	"fmt"
	"reflect"

	"github.com/knyazev-ro/vulcan-orm/orm/model"
)

type GenerateColStringStruct struct {
	TableName string
	ColTag    string
	AggTag    string
}

type GenerateColsOptions struct {
	useAggs bool
}

func (q *Query[T]) generateColString(data *GenerateColStringStruct) (string, string, string) {
	colTag := "all"
	if data.ColTag != "*" {
		colTag = data.ColTag
	}
	original := fmt.Sprintf(`"%s"."%s"`, data.TableName, data.ColTag)

	if data.ColTag == "*" {
		original = "*"
	}

	as := fmt.Sprintf(`%s_%s`, data.TableName, colTag)
	alias := fmt.Sprintf(`%s AS %s`, original, as)
	switch data.AggTag {
	case "sum":
		alias = fmt.Sprintf(`SUM(%s) AS %s_sum`, original, as)
	case "avg":
		alias = fmt.Sprintf(`AVG(%s) AS %s_avg`, original, as)
	case "max":
		alias = fmt.Sprintf(`MAX(%s) AS %s_max`, original, as)
	case "min":
		alias = fmt.Sprintf(`MIN(%s) AS %s_min`, original, as)
	case "count":
		alias = fmt.Sprintf(`COUNT(%s) AS %s_count`, original, as)
	default:
	}
	return alias, as, original
}

func (q *Query[T]) generateCols(i interface{}, options *GenerateColsOptions) []string {
	cols := []string{}
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		panic("Must be a struct")
	}
	metadata, ok := val.Type().FieldByName("_")
	if !ok {
		panic("metadata is not found")
	}
	TableName := metadata.Tag.Get("table")

	originalColsForAgg := []string{}
	shouldGroup := false

	for i := range val.NumField() {
		valueType := val.Type().Field(i)
		typeTag := valueType.Tag.Get("type")
		if typeTag == "column" {
			aggTag := valueType.Tag.Get("agg")
			colTag := valueType.Tag.Get("col")
			tableTag := valueType.Tag.Get("table")

			if !options.useAggs {
				aggTag = ""
			}

			if tableTag == "" {

				if aggTag != "" {
					shouldGroup = true
				} else {
					originalColsForAgg = append(originalColsForAgg, fmt.Sprintf(`%s.%s`, TableName, colTag))
				}
				alias, _, _ := q.generateColString(&GenerateColStringStruct{TableName: TableName, ColTag: colTag, AggTag: aggTag})
				cols = append(cols, alias)
			} else {
				if aggTag != "" {
					shouldGroup = true
				} else {
					originalColsForAgg = append(originalColsForAgg, fmt.Sprintf(`%s.%s`, tableTag, colTag))
				}
				alias, _, _ := q.generateColString(&GenerateColStringStruct{TableName: tableTag, ColTag: colTag, AggTag: aggTag})
				cols = append(cols, alias)

			}
		}
	}

	if shouldGroup && len(originalColsForAgg) > 0 {
		q.GroupBy(originalColsForAgg)
	}

	return cols
}

func (q *Query[T]) MSelect(i interface{}) *Query[T] {
	cols := q.generateCols(i, &GenerateColsOptions{useAggs: true})
	metadata, ok := reflect.TypeOf(i).Elem().FieldByName("_")
	if !ok {
		panic("metadata is not found")
	}
	q.Model = model.Model{
		TableName: metadata.Tag.Get("table"),
		Pk:        metadata.Tag.Get("pk"),
	}
	if len(cols) > 0 {
		q.selectRaw(cols)
	}
	return q
}
