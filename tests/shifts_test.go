package tests

import (
	"testing"
	"fmt"
	)


// Create tests
func Test_HTTP_CreateShiftAsEmployee(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=1", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 0,
		    "start_time": "Thu, Aug 1 19:31:46.631 2018",
		    "end_time": "Thu, Aug 1 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 403 {
			t.Fatal("Error, request should have failed.")
		}
	}
}

func Test_HTTP_CreateShiftValid(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 0,
		    "start_time": "Thu, Aug 1 19:31:46.631 2018",
		    "end_time": "Thu, Aug 1 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 200 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_CreateShiftValid100(t *testing.T) {
	for i := 0; i < 100; i ++ {
		if _, code, err := PostURL("shifts?current_user_id=3", fmt.Sprintf(`
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 0,
		    "start_time": "Thu, Aug 1 19:31:46.631 30%d",
		    "end_time": "Thu, Aug 1 20:31:46.631 30%d"
		}
	`, i, i)); err != nil {
			t.Fatal(err)
		} else {
			//Add unmarshal test.
			if *code != 200 {
				t.Fatal("Error, request should have succeeded.")
			}
		}
	}

}

func Test_HTTP_CreateShiftValidNullEmployee(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": null,
		    "break": 0,
		    "start_time": "Thu, Aug 1 19:31:46.631 2018",
		    "end_time": "Thu, Aug 1 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 200 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_CreateShiftValidNullManager(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "employee_id": null,
		    "break": 0,
		    "start_time": "Thu, Aug 7 19:31:46.631 2018",
		    "end_time": "Thu, Aug 7 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 200 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_CreateShiftOverlapping(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 0,
		    "start_time": "Thu, Aug 2 19:31:46.631 2018",
		    "end_time": "Thu, Aug 2 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 200 {
			t.Fatal("Error, request should have succeeded.")
		}
	}

	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 0,
		    "start_time": "Thu, Aug 2 19:31:46.631 2018",
		    "end_time": "Thu, Aug 2 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have failed.")
		}
	}
}

func Test_HTTP_CreateShiftBadStartEnd(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 0,
		    "start_time": "Thu, Aug 2 19:31:46.631 2018",
		    "end_time": "Thu, Aug 1 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_CreateShiftNullStart(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 0,
		    "start_time": null,
		    "end_time": "Thu, Aug 1 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_CreateShiftWhitespaceStart(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 0,
		    "start_time": "   ",
		    "end_time": "Thu, Aug 1 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_CreateShiftNullEnd(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 0,
		    "start_time": "Thu, Aug 1 20:31:46.631 2018",
		    "end_time": null
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_CreateShiftWhitespaceEnd(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 0,
		    "start_time": "Thu, Aug 1 20:31:46.631 2018",
		    "end_time": "       "
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_CreateShiftLongBreak(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 1,
		    "break": 1,
		    "start_time": "Thu, Aug 3 19:31:46.631 2018",
		    "end_time": "Thu, Aug 3 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 500 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_CreateShiftBadManagerID(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 0,
		    "employee_id": 1,
		    "break": 1,
		    "start_time": "Thu, Aug 5 19:31:46.631 2018",
		    "end_time": "Thu, Aug 5 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 500 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_CreateShiftBadEmployeeID(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		{
		    "id": 3,
		    "manager_id": 3,
		    "employee_id": 0,
		    "break": 1,
		    "start_time": "Thu, Aug 6 19:31:46.631 2018",
		    "end_time": "Thu, Aug 6 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 404 {
			t.Fatal("Error, request should have failed.")
		}
	}
}

func Test_HTTP_CreateShiftMalformedJSON(t *testing.T) {
	if _, code, err := PostURL("shifts?current_user_id=3", `
		
		    "break": 1,
		    "start_time": "Thu, Aug 6 19:31:46.631 2018",
		    "end_time": "Thu, Aug 6 20:31:46.631 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}


// Update tests
func Test_HTTP_UpdateShiftAsManager(t *testing.T) {
	if _, code, err := PutURL("shifts/1?current_user_id=3", `
		{
		    "break": 1,
		    "start_time": "Sun, Aug 19 18:30:00.000 2018",
		    "end_time": "Mon, Aug 19 20:30:00.00 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 200 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_UpdateShiftAsEmployee(t *testing.T) {
	if _, code, err := PutURL("shifts/1?current_user_id=1", `
		{
		    "break": 1,
		    "start_time": "Sun, Aug 19 18:30:00.000 2018",
		    "end_time": "Mon, Aug 19 20:30:00.00 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 403 {
			t.Fatal("Error, request should have failed.")
		}
	}
}

func Test_HTTP_UpdateShiftInvalidStartEnd(t *testing.T) {
	if _, code, err := PutURL("shifts/1?current_user_id=3", `
		{
		    "break": 1,
		    "start_time": "Mon, Aug 19 20:30:00.00 2018",
		    "end_time": "Sun, Aug 19 18:30:00.000 2018"
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have failed.")
		}
	}
}

func Test_HTTP_UpdateShiftSetBreak(t *testing.T) {
	if _, code, err := PutURL("shifts/1?current_user_id=3", `
		{
		    "break": 0.45
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		if *code != 200 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_UpdateShiftSetBreakNegative(t *testing.T) {
	if _, code, err := PutURL("shifts/1?current_user_id=3", `
		{
		    "break": -0.45
		}
	`); err != nil {
		t.Fatal(err)
	} else {
		if *code != 400 {
			t.Fatal("Error, request should have failed.")
		}
	}
}


