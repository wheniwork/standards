package data

import (
	"github.com/ECourant/standards/filtering"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/ECourant/standards/conf"
)

var (
	SummaryConstraints = filtering.GenerateConstraints(Summary{})
)

type DSummary struct {
	DSession
}

func (ctx DSummary) Constraints() filtering.RequestConstraints {
	return SummaryConstraints
}

type Summary struct {
	EmployeeID                  *int     `json:"employee_id,omitempty" query:"15" name:"ID"`
	EmployeeUserObj             *User    `json:"employee_user,omitempty" query:"8" name:"Employee User"`
	WeekStart                   *string  `json:"week_start" query:"11" name:"Week Start" range:"starting"`
	WeekEnd                     *string  `json:"week_end" query:"11" name:"Week End" range:"ending"`
	TotalShifts                 *int     `json:"total_shifts" query:"11" name:"Total Shifts"`
	TotalScheduledTime          *float64 `json:"total_scheduled_time" query:"11" name:"Total Scheduled Time"`
	TotalScheduledTimeFormatted *string  `json:"total_scheduled_time_formatted" query:"11" name:"Total Scheduled Time Formatted"`
	TotalWorkedTime             *float64 `json:"total_worked_time" query:"11" name:"Total Worked Time"`
	TotalWorkedTimeFormatted    *string  `json:"total_worked_time_formatted" query:"11" name:"Total Worked Time Formatted"`
	TotalBreakTime              *float64 `json:"total_break_time" query:"11" name:"Total Break Time"`
	TotalBreakTimeFormatted     *string  `json:"total_break_time_formatted" query:"11" name:"Total Break Time Formatted"`
}

type summaryRow struct {
	Summary
	EmployeeUser *string `json:"employee_user,omitempty" query:"11" name:"Employee User"`
}

func rowToSummary(rows []summaryRow) []Summary {
	result := make([]Summary, len(rows))
	for i, row := range rows {
		Summary := row.Summary
		if row.EmployeeUser != nil {
			json.Unmarshal([]byte(*row.EmployeeUser), &Summary.EmployeeUserObj)
		}
		result[i] = Summary
	}
	return result
}

func (ctx DSummary) GetSummary(id *int, params filtering.RequestParams) ([]Summary, *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	db.LogMode(true)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve summary at this time.", err)
	}
	defer db.Close()

	result := make([]summaryRow, 0)

	db = db.
		Table("public.vw_shifts_summary_api").
		Select(params.Fields).
		Order(params.Sorts).
		Offset((params.Page * params.PageSize) - params.PageSize).
		Limit(params.PageSize)
	if id != nil {
		db = db.Where("employee_id = ?", id)
	}

	if len(params.Filters) > 0 || params.DateRange != nil {
		db = filtering.WhereFilters(db, params, ctx.Constraints())
	}

	db.Scan(&result)
	return rowToSummary(result), nil
}
