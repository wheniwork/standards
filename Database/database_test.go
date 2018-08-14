package Database

import (
	"testing"
		_ "github.com/lib/pq"
	"database/sql"
	"os"
	"github.com/ecourant/standards/Site/conf"
	)

var (
	Config conf.Config
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


func Test_CreateUser(t *testing.T) {
	if err := runQueryWithRollback(t, "INSERT INTO public.users (name, email, phone, role) VALUES($1, $2, $3, $4);", "Billy", "billy@mays.com", nil, "employee"); err != nil {
		t.Error(err)
		t.Fail()
	}
}

func Test_CreateUserInvalidRole(t *testing.T) {
	if err := runQueryWithRollback(t, "INSERT INTO public.users (name, email, phone, role) VALUES($1, $2, $3, $4);", "Billy", "billy@mays.com", nil, "admin"); err == nil {
		t.Errorf("insert should have failed because of invalid role 'admin'")
		t.Fail()
	}
}

func Test_CreateShiftForUser(t *testing.T) {
	if err := runQueryWithRollback(t, "INSERT INTO public.shifts (manager_id,employee_id,start_time,end_time) VALUES($1, $2, $3::timestamp, $4::timestamp);", 3, 1, "2018-08-13 8:00AM", "2018-08-13 4:00PM" ); err != nil {
		t.Error(err)
		t.Fail()
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
