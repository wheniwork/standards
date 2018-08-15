package filtering

import (
	"reflect"
	"fmt"
	"github.com/kataras/iris"
	"encoding/json"
	"strings"
	"strconv"
	"github.com/jinzhu/gorm"
	"github.com/kataras/go-errors"
)

type Filter struct {
	Field      string    `json:"field,omitempty"`
	Equals     *string   `json:"equals,omitempty"`
	EqualsInt  *int      `json:"equals_int,omitempty"`
	EqualsBool *bool     `json:"equals_bool,omitempty"`
	Like       *string   `json:"like,omitempty"`
	NotLike    *string   `json:"not_like,omitempty"`
	InInt      *[]int64  `json:"in_int,omitempty"`
	In         *[]string `json:"in,omitempty"`
}

type Sort struct {
	Field     string
	Direction string
}

type RequestParams struct {
	Page      int
	PageSize  int
	Fields    string
	Sorts     string
	Filters   []Filter
	DateRange *DateRange
}

type DateRange struct {
	Start *string
	End   *string
}

type RequestConstraints struct {
	Fields             map[string]RequestQueryField
	QueryTypeCount     map[QueryType]int
	DefaultSort        string
	StartingRangeField *string
	EndingRangeField   *string
}

type RequestQueryField struct {
	QueryType QueryType
	Name      string
	Type      reflect.Type
}

type QueryType int
type RequestType int

const (
	Filterable      QueryType = 1 << iota
	Sortable        QueryType = 2
	DefaultSortAsc  QueryType = 4
	DefaultSortDesc QueryType = 16
	Standard        QueryType = 8 // Next should be 32
)

const (
	StandardRequest RequestType = 1 << iota
	DetailedRequest RequestType = 2
	CountRequest    RequestType = 4
)

var (
	QueryTypeNames = map[string]QueryType{

		"Fiterable":   1,
		"Sortable":    2,
		"DefaultSort": 4,
		"Standard":    8,
	}
)

func GenerateConstraints(T interface{}) RequestConstraints {
	cons := RequestConstraints{
		Fields:         map[string]RequestQueryField{},
		QueryTypeCount: map[QueryType]int{},
	}
	defaultSorts := make([]string, 0)
	ref := reflect.TypeOf(T)
	for i := 0; i < ref.NumField(); i++ {
		f := ref.Field(i)
		fieldName := strings.Split(f.Tag.Get("json"), ",")[0]
		if val, ok := f.Tag.Lookup("query"); !ok {
			fmt.Println("ALERT: Field (", fieldName, ") has no constraint parameters and will not have functionality but may still be returned.")
		} else if q, err := strconv.Atoi(val); err != nil {
			panic("Error, field (" + fieldName + ") failed to be parsed as an int.")
		} else {
			if name, ok := f.Tag.Lookup("name"); ok {
				cons.Fields[fieldName] = RequestQueryField{
					QueryType: QueryType(q),
					Name:      name,
					Type:      f.Type,
				}
			} else {
				cons.Fields[fieldName] = RequestQueryField{
					QueryType: QueryType(q),
					Type:      f.Type,
				}
			}

			for _, n := range QueryTypeNames {
				if QueryType(q)|n == QueryType(q) {
					cons.QueryTypeCount[n]++
				}
			}

			if QueryType(q)|DefaultSortAsc == QueryType(q) {
				defaultSorts = append(defaultSorts, fieldName+" ASC")
			} else if QueryType(q)|DefaultSortDesc == QueryType(q) {
				defaultSorts = append(defaultSorts, fieldName+" DESC")
			}
		}

		if val, ok := f.Tag.Lookup("range"); ok {
			switch val {
			case "starting":
				cons.StartingRangeField = &fieldName
			case "ending":
				cons.EndingRangeField = &fieldName
			}
		}
	}
	if len(defaultSorts) > 0 {
		cons.DefaultSort = strings.Join(defaultSorts, ",")
	}
	return cons
}

func ParseRequestParams(ctx iris.Context, constraints RequestConstraints, requestType RequestType) (*RequestParams, error) {
	params := RequestParams{}
	if requestType|StandardRequest == requestType {
		params.Page = ctx.URLParamIntDefault("page", 1)
		params.PageSize = ctx.URLParamIntDefault("page_size", 10)

		if sorts := strings.Split(ctx.URLParam("order"), ","); len(sorts) > 0 && sorts[0] != "" {
			if val, ok := constraints.QueryTypeCount[Sortable]; !(ok && val > 0) {
				return nil, errors.New("Error, sorting is not allowed for this request.")
			}
			st := make([]string, 0)
			dirs := map[byte]string{
				byte('a'): "ASC",
				byte('d'): "DESC",
			}
			for _, s := range sorts {
				if direction, ok := dirs[byte(s[0])]; ok {
					if val, ok := constraints.Fields[string(s[1:])]; ok && val.QueryType|Sortable == val.QueryType && val.QueryType|Standard == val.QueryType {
						st = append(st, string(s[1:])+" "+direction)
					} else {
						return nil, errors.New("Error, could not sort by field (" + string(s[1:]) + "), field is not sortable.")
					}
				} else {
					return nil, errors.New("Error, cannot parse sort parameter (" + s + "), direction not valid.")
				}
			}
			if len(st) == 0 {
				return nil, errors.New("Error, no valid sortable fields were specified.")
			} else {
				params.Sorts = strings.Join(st, ",")
			}
		} else {
			params.Sorts = constraints.DefaultSort
		}
	}

	if requestType|StandardRequest == requestType || requestType|DetailedRequest == requestType {
		if fields := strings.Split(ctx.URLParam("fields"), ","); len(fields) > 0 && fields[0] != "" {
			ft := make([]string, 0)
			for _, f := range fields {
				if fd, ok := constraints.Fields[f]; ok {
					if requestType|StandardRequest == requestType && !(fd.QueryType|Standard == fd.QueryType) {
						return nil, errors.New("Error, cannot select field (" + f + ") it is only valid for detailed requests.")
					} else {
						ft = append(ft, f)
					}
				} else {
					return nil, errors.New("Error, cannot select field (" + f + ") it is not valid.")
				}
			}
			if len(ft) == 0 {
				return nil, errors.New("Error, no valid fields were specified.")
			} else {
				params.Fields = strings.Join(ft, ",")
			}
		} else {
			params.Fields = "*"
		}
	}

	if requestType|StandardRequest == requestType || requestType|DetailedRequest == requestType || requestType|CountRequest == requestType {
		if f := ctx.URLParam("filter"); f != "" {
			if val, ok := constraints.QueryTypeCount[Filterable]; !(ok && val > 0) {
				return nil, errors.New("Error, filtering is not allowed for this request.")
			}
			filters := make([]Filter, 0)
			if err := json.Unmarshal([]byte(f), &filters); err != nil {
				return nil, errors.New("Error, could not parse filter.")
			} else {
				params.Filters = make([]Filter, 0)
				for _, filter := range filters {
					if val, ok := constraints.Fields[filter.Field]; !(ok && val.QueryType|Filterable == val.QueryType) {
						return nil, errors.New("Error, cannot filter by field (" + filter.Field + "), it is not valid for this request.")
					} else {
						params.Filters = append(params.Filters, filter)
					}
				}
			}
		} else {
			params.Filters = make([]Filter, 0)
		}
		start, end := ctx.URLParam("date_from"), ctx.URLParam("date_to")
		daterange := DateRange{}
		if start != "" {
			daterange.Start = &start
		}
		if end != "" {
			daterange.End = &end
		}
		params.DateRange = &daterange
	}

	return &params, nil
}

func WhereFilters(db *gorm.DB, params RequestParams, constraints RequestConstraints) *gorm.DB {
	for _, filter := range params.Filters {
		t := struct {
			st *string
			it *int
			bo *bool
		}{}
		if field, ok := constraints.Fields[filter.Field]; ok && field.Type == reflect.TypeOf(t.st) {
			if filter.Equals != nil {
				db = db.Where(filter.Field+" = ?", filter.Equals)
			} else if filter.In != nil {
				db = db.Where(filter.Field+" IN (?)", filter.In)
			} else if filter.Like != nil {
				db = db.Where(filter.Field+" LIKE ?", filter.Like)
			} else if filter.NotLike != nil {
				db = db.Where(filter.Field+" NOT LIKE ?", filter.NotLike)
			}
		} else if ok && field.Type == reflect.TypeOf(t.it) {
			if filter.EqualsInt != nil {
				db = db.Where(filter.Field+" = ?", filter.EqualsInt)
			} else if filter.InInt != nil {
				db = db.Where(filter.Field+" IN (?)", filter.InInt)
			}
		} else if ok && field.Type == reflect.TypeOf(t.bo) {
			if filter.EqualsBool != nil {
				db = db.Where(filter.Field+" = ?", filter.EqualsBool)
			}
		} else {
			fmt.Println("Type", field.Type)
		}
	}
	if constraints.StartingRangeField != nil && constraints.EndingRangeField != nil {
		if params.DateRange != nil && (params.DateRange.Start != nil || params.DateRange.End != nil) {
			if params.DateRange.Start != nil {
				db = db.Where(fmt.Sprintf("%s::date >= ?::date", *constraints.StartingRangeField), *params.DateRange.Start)
			}

			if params.DateRange.End != nil {
				db = db.Where(fmt.Sprintf("%s::date <= ?::date", *constraints.EndingRangeField), *params.DateRange.End)
			}
		}
	}

	return db
}

func GetFilter(ctx iris.Context) []Filter {
	if f := ctx.URLParam("filter"); f != "" {
		filters := make([]Filter, 0)
		json.Unmarshal([]byte(f), &filters)
		return filters
	} else {
		return make([]Filter, 0)
	}
}
