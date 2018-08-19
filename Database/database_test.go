package Database

import (
	"testing"
	_ "github.com/lib/pq"
	"database/sql"
	"os"
	"github.com/ecourant/standards/Site/conf"
)

var (
	Config  conf.Config
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

func TestMain(m *testing.M) {
	if conf, err := conf.LoadConfig("database_test_config.json"); err != nil {
		panic(err)
	} else {
		Config = *conf
	}

	retCode := m.Run()
	os.Exit(retCode)
}

// Just run all of the queries in the array and verify whether or not the response was supposed to be successful.
func Test_HTTP_Inserts(t *testing.T) {
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
	db, err := sql.Open("postgres", Config.ConnectionString)
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
