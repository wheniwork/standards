package tests

import (
	"testing"
)

func Test_HTTP_GetUsersWithoutCurrentID(t *testing.T) {
	if _, code, err := GetURL("users"); err != nil {
		t.Fatal(err)
	} else {
		if *code != 403 {
			t.Fatal("Error, request should fail as unauthorized.")
		}
	}
}

func Test_HTTP_GetUsersWithValidCurrentID(t *testing.T) {
	if _, code, err := GetURL("users?current_user_id=1"); err != nil {
		t.Fatal(err)
	} else {
		// Add unmarshal test.
		if *code != 200 {
			t.Fatal("Error, request should have succeeded.")
		}
	}
}

func Test_HTTP_GetUsersWithNegativeCurrentID(t *testing.T) {
	if _, code, err := GetURL("users?current_user_id=-1"); err != nil {
		t.Fatal(err)
	} else {
		// Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have failed with client error.")
		}
	}
}

func Test_HTTP_GetUsersWithNonIntegerCurrentID(t *testing.T) {
	if _, code, err := GetURL("users?current_user_id=vdassda"); err != nil {
		t.Fatal(err)
	} else {
		// Add unmarshal test.
		if *code != 400 {
			t.Fatal("Error, request should have failed with client error.")
		}
	}
}
