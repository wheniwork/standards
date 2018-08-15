package data

import (
	"github.com/ecourant/standards/Site/filtering"
	"github.com/jinzhu/gorm"
	"github.com/ecourant/standards/Site/conf"
)

var (
	ShiftConstraints = filtering.GenerateConstraints(Shift{})
)

type DShifts struct {
	DSession
}

func (ctx DShifts) Constraints() filtering.RequestConstraints {
	return ShiftConstraints
}

type Shift struct {
	ID         *int     `json:"id,omitempty" query:"15" name:"ID"`
	ManagerID  *int     `json:"manager_id,omitempty" query:"11" name:"Manager ID"`
	EmployeeID *int     `json:"employee_id,omitempty" query:"11" name:"Employee ID"`
	Break      *float64 `json:"break,omitempty" query:"11" name:"Break"`
	StartTime  *string  `json:"start_time,omitempty" query:"11" name:"Start Time" range:"starting"`
	EndTime   *string  `json:"end_time,omitempty" query:"11" name:"End Time" range:"ending"`
	CreatedAt  *string  `json:"created_at,omitempty" query:"11" name:"Created At"`
	UpdatedAt  *string  `json:"updated_at,omitempty" query:"11" name:"Updated At"`
}

func (ctx DShifts) GetShifts(params filtering.RequestParams) ([]Shift, *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	db.LogMode(true)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve shifts at this time.", err)
	}
	defer db.Close()

	result := make([]Shift, 0)

	db = db.
		Table("public.vw_shifts_api").
		Select(params.Fields).
		Order(params.Sorts).
		Offset((params.Page * params.PageSize) - params.PageSize).
		Limit(params.PageSize)

	if len(params.Filters) > 0 || params.DateRange != nil {
		db = filtering.WhereFilters(db, params, ctx.Constraints())
	}

	db.Scan(&result)
	return result, nil
}

func (ctx DShifts) GetMyShifts(params filtering.RequestParams) ([]Shift, *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	db.LogMode(true)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve shifts at this time.", err)
	}
	defer db.Close()

	result := make([]Shift, 0)

	db = db.
		Table("public.vw_shifts_api").
		Select(params.Fields).
		Order(params.Sorts).
		Offset((params.Page * params.PageSize) - params.PageSize).
		Limit(params.PageSize).
		Where("(employee_id = ? OR employee_id IS NULL)", ctx.UserID)

	if len(params.Filters) > 0 || params.DateRange != nil {
		db = filtering.WhereFilters(db, params, ctx.Constraints())
	}

	db.Scan(&result)
	return result, nil
}

func (ctx DShifts) GetMySummary(params filtering.RequestParams) ([]Shift, *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve shifts at this time.", err)
	}
	defer db.Close()

	result := make([]Shift, 0)

	db = db.
		Table("public.vw_shifts_api").
		Select(params.Fields).
		Order(params.Sorts).
		Offset((params.Page * params.PageSize) - params.PageSize).
		Limit(params.PageSize).Where("(employee_id = ? OR manager_id = ?)", ctx.UserID, ctx.UserID)

	if len(params.Filters) > 0 {
		db = filtering.WhereFilters(db, params, ctx.Constraints())
	}

	db.Scan(&result)
	return result, nil
}

func (ctx DShifts) GetMyShiftDetails(params filtering.RequestParams) ([]Shift, *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve shifts at this time.", err)
	}
	defer db.Close()

	result := make([]Shift, 0)

	db = db.
		Table("public.vw_shifts_api").
		Select(params.Fields).
		Order(params.Sorts).
		Offset((params.Page * params.PageSize) - params.PageSize).
		Limit(params.PageSize).Where("(employee_id = ? OR manager_id = ?)", ctx.UserID, ctx.UserID)

	if len(params.Filters) > 0 {
		db = filtering.WhereFilters(db, params, ctx.Constraints())
	}

	db.Scan(&result)
	return result, nil
}