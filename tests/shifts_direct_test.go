package tests

import (
	"testing"
	"github.com/ECourant/standards/data"
)

func getContext() data.DSession {
	return data.DSession{UserID: 3, IsManager: true}
}

func Test_CreateShift(t *testing.T) {
	EmployeeID, ManagerID := 3, 3
	StartTime, EndTime := "Thu, Aug 1 20:00:00.00 2018", "Thu, Aug 1 20:31:46.631 2018"

	if result, err := getContext().Shifts().CreateShift(data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have failed to create shift as employee.")
		}
	}
}

func Test_CreateShiftBadTime(t *testing.T) {
	EmployeeID, ManagerID := 3, 3
	EndTime, StartTime := "Thu, Aug 1 20:00:00.00 2018", "Thu, Aug 1 20:31:46.631 2018"

	if result, err := getContext().Shifts().CreateShift(data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {

	} else {
		if result != nil {
			t.Fatal("Error, should have failed to create shift as employee.")
		}
	}
}

func Test_CreateShiftUpdateEmployee(t *testing.T) {
	EmployeeID, ManagerID := 3, 3
	StartTime, EndTime := "01-01-2018 8:00AM", "01-01-2018 4:00PM"
	id := 0
	if result, err := getContext().Shifts().CreateShift(data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have failed to create shift as employee.")
		} else {
			id = *result.ID
		}
	}
	EmployeeID = 2
	if result, err := getContext().Shifts().UpdateShift(id, data.Shift{
		EmployeeID: &EmployeeID,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have failed to create shift as employee.")
		} else if *result.EmployeeID != EmployeeID {
			t.Fatal("Error, updated employee ID does not match.")
		}
	}
}

func Test_UpdateShiftConflicting(t *testing.T) {
	EmployeeID, ManagerID := 3, 3
	StartTime, EndTime := "02-02-2018 8:00AM", "02-02-2018 10:00AM"
	id := 0
	if result, err := getContext().Shifts().CreateShift(data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have failed to create shift as employee.")
		} else {
			id = *result.ID
		}
	}
	StartTime, EndTime = "02-02-2018 10:00AM", "02-02-2018 11:00AM"
	if result, err := getContext().Shifts().CreateShift(data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have failed to create shift as employee.")
		} else {
			id = *result.ID
		}
	}
	StartTime, EndTime = "02-02-2018 9:00AM", "02-02-2018 11:00AM"
	if result, err := getContext().Shifts().UpdateShift(id, data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {

	} else {
		if result != nil {
			t.Fatal("Error, should have failed to create shift as employee.")
		}
	}
}

func Test_CreateShiftConflicting(t *testing.T) {
	EmployeeID, ManagerID := 3, 3
	StartTime, EndTime := "02-02-2018 8:00AM", "02-02-2018 10:00AM"
	if result, err := getContext().Shifts().CreateShift(data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have failed to create shift as employee.")
		} else {

		}
	}
	StartTime, EndTime = "02-02-2018 7:00AM", "02-02-2018 9:00AM"
	if _, err := getContext().Shifts().CreateShift(data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {

	} else {
		t.Fatal("Error, should have failed to create shift as employee.")
	}
}

func Test_UpdateShiftRemoveEmployee(t *testing.T) {
	EmployeeID, ManagerID := 3, 3
	StartTime, EndTime := "03-03-2018 8:00AM", "03-03-2018 10:00AM"
	id := 0
	if result, err := getContext().Shifts().CreateShift(data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have failed to create shift.")
		} else {
			id = *result.ID
		}
	}
	EmployeeID = -1
	if result, err := getContext().Shifts().UpdateShift(id, data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have succeeded in updating shift..")
		} else if result.EmployeeID != nil {
			t.Fatal("Error, employee should be nil.")
		}
	}
}

func Test_UpdateShiftAddBreak(t *testing.T) {
	EmployeeID, ManagerID := 3, 3
	StartTime, EndTime := "03-04-2018 8:00AM", "03-04-2018 10:00AM"
	id := 0
	if result, err := getContext().Shifts().CreateShift(data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have failed to create shift.")
		} else {
			id = *result.ID
		}
	}
	Break := 0.5
	if result, err := getContext().Shifts().UpdateShift(id, data.Shift{
		Break: &Break,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have succeeded in updating shift..")
		} else if result.Break != nil && *result.Break != Break {
			t.Fatal("Error, break should be updated.")
		}
	}
}

func Test_DeleteShift(t *testing.T) {
	EmployeeID, ManagerID := 3, 3
	StartTime, EndTime := "03-05-2018 8:00AM", "03-05-2018 10:00AM"
	id := 0
	if result, err := getContext().Shifts().CreateShift(data.Shift{
		EmployeeID: &EmployeeID,
		ManagerID:  &ManagerID,
		StartTime:  &StartTime,
		EndTime:    &EndTime,
	}); err != nil {
		t.Fatal(err)
	} else {
		if result == nil {
			t.Fatal("Error, should have failed to create shift.")
		} else {
			id = *result.ID
		}
	}
	if err := getContext().Shifts().DeleteShift(id); err != nil {
		t.Fatal(err)
	}
}
