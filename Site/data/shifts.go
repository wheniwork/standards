package data

import (
	"github.com/ecourant/standards/Site/filtering"
	"github.com/jinzhu/gorm"
	"github.com/ecourant/standards/Site/conf"
	"encoding/json"
	"fmt"
	"strings"
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
	ID              *int     `json:"id,omitempty" query:"27" name:"ID"`
	ManagerID       *int     `json:"manager_id,omitempty" query:"11" name:"Manager ID"`
	ManagerUserObj  *User    `json:"manager_user,omitempty" query:"8" name:"Manager User"`
	EmployeeID      *int     `json:"employee_id,omitempty" query:"11" name:"Employee ID"`
	EmployeeUserObj *User    `json:"employee_user,omitempty" query:"8" name:"Employee User"`
	Break           *float64 `json:"break,omitempty" query:"11" name:"Break"`
	StartTime       *string  `json:"start_time,omitempty" query:"11" name:"Start Time" range:"starting"`
	EndTime         *string  `json:"end_time,omitempty" query:"11" name:"End Time" range:"ending"`
	CreatedAt       *string  `json:"created_at,omitempty" query:"11" name:"Created At"`
	UpdatedAt       *string  `json:"updated_at,omitempty" query:"11" name:"Updated At"`
}

type shiftRow struct {
	Shift
	ManagerUser  *string `json:"manager_user,omitempty" query:"11" name:"Manager User"`
	EmployeeUser *string `json:"employee_user,omitempty" query:"11" name:"Employee User"`
}


func rowsToShifts(rows []shiftRow) []Shift {
	result := make([]Shift, len(rows))
	for i, row := range rows {
		shift := row.Shift
		if row.ManagerUser != nil {
			json.Unmarshal([]byte(*row.ManagerUser), &shift.ManagerUserObj)
		}
		if row.EmployeeUser != nil {
			json.Unmarshal([]byte(*row.EmployeeUser), &shift.EmployeeUserObj)
		}
		result[i] = shift
	}
	return result
}

func (ctx DShifts) GetShifts(params filtering.RequestParams) ([]Shift, *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	db.LogMode(true)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve shifts at this time.", err)
	}
	defer db.Close()

	result := make([]shiftRow, 0)

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
	return rowsToShifts(result), nil
}

func (ctx DShifts) GetMyShifts(params filtering.RequestParams) ([]Shift, *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	db.LogMode(true)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve shifts at this time.", err)
	}
	defer db.Close()

	result := make([]shiftRow, 0)

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
	return rowsToShifts(result), nil
}

func (ctx DShifts) GetShiftDetails(params filtering.RequestParams, id int) ([]Shift, *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	db.LogMode(true)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve shifts at this time.", err)
	}
	defer db.Close()

	result := make([]shiftRow, 0)

	db = db.
		Table("public.vw_shifts_detailed_api").
		Select(params.Fields).
		Order(params.Sorts).
		Offset((params.Page * params.PageSize) - params.PageSize).
		Limit(params.PageSize).
		Where("group_by_id = ?", id)

	if len(params.Filters) > 0 || params.DateRange != nil {
		db = filtering.WhereFilters(db, params, ctx.Constraints())
	}

	db.Scan(&result)
	return rowsToShifts(result), nil
}

func (ctx DShifts) CreateShift(shift Shift) (response *Shift, rerr *DError) {

	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	db.LogMode(true)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve shifts at this time.", err)
	}
	defer db.Close()

	result := make([]shiftRow, 0)

	db = db.Begin()
	// I hate how this looks in golang, but basically if there is a panic or an error somewhere further down, the transaction will rollback.
	defer func() {
		if r := recover(); r != nil {
			db.Rollback()
			response = nil
			rerr = NewServerError("Error, could not create shift at this time.", err)
			return
		}
	}()

	if err := ctx.verifyShift(nil, &shift, db); err != nil {
		return nil, err
	}

	if err := db.Raw(`INSERT INTO public.shifts (manager_id,employee_id,break,start_time,end_time)
				 VALUES(?, ?, ?, ?::timestamp, ?::timestamp) 
        		 RETURNING 
					id,
					manager_id,
					employee_id,
					break,
					to_char(start_time, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS start_time,
					to_char(end_time, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS end_time,
					to_char(created_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS created_at,
					to_char(updated_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS updated_at;`, shift.ManagerID, shift.EmployeeID, shift.Break, shift.StartTime, shift.EndTime).Scan(&result).Error; err != nil {
		db.Rollback()
		return nil, NewServerError("Error, an unexpected error occurred. The shift was not created.", err)
	}
	db.Commit()
	return &rowsToShifts(result)[0], nil
}

func (ctx DShifts) UpdateShift(id int, shift Shift) (response *Shift, rerr *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	db.LogMode(true)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve shifts at this time.", err)
	}
	defer db.Close()

	db = db.Begin()
	// I hate how this looks in golang, but basically if there is a panic or an error somewhere further down, the transaction will rollback.
	defer func() {
		if r := recover(); r != nil {
			db.Rollback()
			response = nil
			rerr = NewServerError("Error, could not update shift at this time.", err)
			return
		}
	}()

	if err := ctx.verifyShift(&id, &shift, db); err != nil {
		return nil, err
	}

	result := make([]shiftRow, 0)

	if err := db.Raw(`
		UPDATE public.shifts SET
			manager_id=COALESCE(?, manager_id),
			employee_id=NULLIF(COALESCE(?, manager_id), -1),
			break=COALESCE(?, break),
			start_time=COALESCE(?::timestamp, start_time),
			end_time=COALESCE(?::timestamp, end_time),
			updated_at=LOCALTIMESTAMP
		WHERE id=?
		RETURNING id,
				  manager_id,
				  employee_id,
				  break,
				  to_char(start_time, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS start_time,
				  to_char(end_time, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS end_time,
				  to_char(created_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS created_at,
				  to_char(updated_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS updated_at;
	`, 	shift.ManagerID, shift.EmployeeID, shift.Break, shift.StartTime, shift.EndTime, id).Scan(&result).Error; err != nil {
		db.Rollback()
		return nil, NewServerError("Error, an unexpected error occurred. The shift was not updated.", err)
	}
	db.Commit()
	return &rowsToShifts(result)[0], nil
}

func (ctx DShifts) DeleteShift(id int) (rerr *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	db.LogMode(true)
	if err != nil {
		return NewServerError("Error, could not delete shift at this time.", err)
	}
	defer db.Close()
	db = db.Begin()
	// I hate how this looks in golang, but basically if there is a panic or an error somewhere further down, the transaction will rollback.
	defer func() {
		if r := recover(); r != nil {
			db.Rollback()
			rerr = NewServerError("Error, could not delete shift at this time.", err)
			return
		}
	}()
	if err := db.Exec("DELETE FROM public.shifts WHERE id = ?").Error; err != nil {
		db.Rollback()
		return NewServerError(fmt.Sprintf("Error, failed to delete shift ID %d.", id), err)
	}
	db.Commit()
	return nil
}

func (ctx DShifts) verifyShift(id *int, shift *Shift , db *gorm.DB) *DError {
	if shift.Break != nil {
		if *shift.Break < 0 {
			return NewClientError("Error, break must be non-negative.", nil)
		}
	}

	if id == nil {
		if shift.StartTime == nil || strings.TrimSpace(*shift.StartTime) == "" {
			return NewClientError("Error, start_time cannot be null or blank.", nil)
		}

		if shift.EndTime == nil || strings.TrimSpace(*shift.EndTime) == "" {
			return NewClientError("Error, end_time cannot be null or blank.", nil)
		}
	}

	if shift.ManagerID == nil {
		shift.ManagerID = &ctx.UserID
	}

	if role, err := ctx.Users().GetUserRole(*shift.ManagerID); err != nil {
		return NewServerError("Error, could not verify manager_id.", err)
	} else if *role == "employee" {
		return NewClientError(fmt.Sprintf("Error, user ID %d is not a manager.", *shift.ManagerID), nil)
	} else if role == nil {
		return NewNotFoundError(fmt.Sprintf("Error, manager_id %d does not exist.", *shift.ManagerID))
	}

	// If the employee id is not null we want to verify
	// that this shift will not conflict with another shift.
	if shift.EmployeeID != nil && *shift.EmployeeID != -1 {
		if role, err := ctx.Users().GetUserRole(*shift.EmployeeID); err != nil {
			return NewServerError("Error, could not verify employee_id.", err)
		}  else if role == nil {
			return NewNotFoundError(fmt.Sprintf("Error, employee_id %d does not exist.", *shift.ManagerID))
		}

		ids := make([]struct {
			ID string
		}, 0)
		d := db.
			Table("public.vw_shifts_api").
			Select("id").
			Where("employee_id = ?", *shift.EmployeeID).
			Where("(start_time::timestamp BETWEEN ?::timestamp AND ?::timestamp) OR (end_time::timestamp BETWEEN ?::timestamp AND ?::timestamp)",
				*shift.StartTime, *shift.EndTime, *shift.StartTime, *shift.EndTime)
		if id != nil { // If this is an update, make sure we exclude the existing shift.
			d = d.Where("id != ?", *id)
		}
		d.Scan(&ids)
		if len(ids) > 0 {
			db.Rollback()
			conflictingShifts := make([]string, len(ids))
			for i, shiftid := range ids {
				conflictingShifts[i] = shiftid.ID
			}
			return NewClientError(fmt.Sprintf("Error, %d shift(s) already exist for user ID %d during the start -> end time. Conflicting shift(s): %s.", len(ids), *shift.EmployeeID, strings.Join(conflictingShifts, ", ")), nil)
		}
	}


	valid_start_end := make([]struct {
		Valid bool
	}, 0)
	if id == nil || (shift.StartTime != nil && shift.EndTime != nil) {
		db.Raw("SELECT ?::timestamp < ?::timestamp AS valid", *shift.StartTime, *shift.EndTime).Scan(&valid_start_end)
		if len(valid_start_end) > 0 {
			if !valid_start_end[0].Valid {
				db.Rollback()
				return NewClientError(fmt.Sprintf("Error, (start_time: %s) must come before (end_time: %s).", *shift.StartTime, *shift.EndTime), nil)
			}
		}
	} else if shift.StartTime != nil || shift.EndTime != nil {
		if shift.StartTime != nil {
			db.Raw("SELECT ?::timestamp < end_time AS valid FROM public.shifts WHERE id = ?;", *shift.StartTime, *id)
		} else if shift.EndTime != nil {
			db.Raw("SELECT ?::timestamp > start_time AS valid FROM public.shifts WHERE id = ?;", *shift.EndTime, *id)
		}
		if len(valid_start_end) > 0 {
			if !valid_start_end[0].Valid {
				db.Rollback()
				if shift.StartTime != nil {
					return NewClientError(fmt.Sprintf("Error, (start_time: %s) must come before end_time.", *shift.StartTime), nil)
				} else {
					return NewClientError(fmt.Sprintf("Error, (end_time: %s) must come after start_time.", *shift.EndTime), nil)
				}
			}
		}
	}

	if id != nil {
		count := 0
		db.
			Table("public.vw_shifts_api").
			Where("id = ?", *id).
			Count(&count)
		if count != 1 {
			return NewNotFoundError(fmt.Sprintf("Error, shift ID %d cannot be updated because it doesn't exist.", *id))
		}
	}
	return nil
}