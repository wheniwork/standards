package tests

import (
	"testing"
	"os"
	"github.com/ecourant/standards/Site/conf"
	"io/ioutil"
	"database/sql"
	"fmt"
)

func TestMain(m *testing.M) {
	if c, err := conf.LoadConfig("test_config.json"); err != nil {
		panic(err)
	} else {
		conf.Cfg = *c
	}
	ResetDatabase()
	StartServer()
	retCode := m.Run()
	os.Exit(retCode)
}

func ResetDatabase() {
	path := "../../Database/Create Database.sql"
	if bytes, err := ioutil.ReadFile(path); err != nil {
		panic(err)
	} else {
		db, err := sql.Open("postgres", conf.Cfg.ConnectionString)
		if err != nil {
			panic(err)
		}
		if tx, err := db.Begin(); err != nil {
			panic(err)
		} else {
			if _, err := tx.Exec(string(bytes)); err != nil {
				panic(err)
			} else {
				tx.Commit()
				fmt.Printf("SUCCESSFULLY RESET DATABASE!\n")
			}
		}
	}
}