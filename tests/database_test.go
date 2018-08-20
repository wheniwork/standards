package tests

import (
	"testing"
	_ "github.com/lib/pq"
	"database/sql"
		"github.com/ECourant/standards/conf"
	)

var (
	Queries = []struct {
		Query      string
		Args       []interface{}
		ShouldFail bool
	}{
		{
			Query: "INSERT INTO public.users (name, email, phone, role) VALUES($1, $2, $3, $4);",
			Args: []interface{}{
				"Billy",
				"billy@mays.com",
				nil,
				"employee",
			},
			ShouldFail: false,
		},
		{
			Query: "INSERT INTO public.users (name, email, phone, role) VALUES($1, $2, $3, $4);",
			Args: []interface{}{
				"Billy",
				"billy@mays.com",
				nil,
				"admin",
			},
			ShouldFail: true,
		},
		{
			Query: "INSERT INTO public.users (name, email, phone, role) VALUES($1, $2, $3, $4);",
			Args: []interface{}{
				"Billy",
				nil,
				nil,
				"admin",
			},
			ShouldFail: true,
		},
		{
			Query: "INSERT INTO public.shifts (manager_id,employee_id,start_time,end_time) VALUES($1, $2, $3::timestamp, $4::timestamp);",
			Args: []interface{}{
				3,
				1,
				"2018-08-13 8:00AM",
				"2018-08-13 4:00PM",
			},
			ShouldFail: false,
		},
	}
)

// Just run all of the queries in the array and verify whether or not the response was supposed to be successful.
func Test_DB_Inserts(t *testing.T) {
	for _, q := range Queries {
		err := runQueryWithRollback(t, q.Query, q.Args...)
		if q.ShouldFail && err == nil {
			t.Errorf("query `%s` should have failed", q.Query)
			t.Fail()
		} else if !q.ShouldFail && err != nil {
			t.Error(err)
			t.Fail()
		}
	}
}

func runQueryWithRollback(t *testing.T, Query string, Args ...interface{}) error {
	db, err := sql.Open("postgres", conf.Cfg.ConnectionString)
	if err != nil {
		t.Error(err)
		t.Fail()
		return err
	}
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		t.Error(err)
		t.Fail()
		return err
	}
	if _, err := tx.Exec(Query, Args...); err != nil {
		return err
	} else {
		return nil
	}
}
