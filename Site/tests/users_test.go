package tests



import (
	"testing"
)


func Test_GetUsersWithoutCurrentID(t *testing.T) {
	if body, code, err := GetURL("users"); err != nil {
		t.Fatal(err)
	} else {
		if *body != `{"message":"Error, current_user_id url param must be specified!","success":false}` {
			t.Fatal("Error, should not be able to retrieve users without current_user_id.")
		}
		if *code != 403 {
			t.Fatal("Error, request should fail as unauthorized.")
		}
	}
}

func Test_GetUsersWithValidCurrentID(t *testing.T) {
	if _, code, err := GetURL("users?current_user_id=1"); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 200 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_GetUsersWithNegativeCurrentID(t *testing.T) {
	if _, code, err := GetURL("users?current_user_id=-1"); err != nil {
		t.Fatal(err)
	} else {
		//Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have failed with client error.")
		}
	}
}