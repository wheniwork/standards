package tests

import (
	"testing"
	"os"
)

func TestMain(m *testing.M) {
	StartServer()
	retCode := m.Run()
	os.Exit(retCode)
}