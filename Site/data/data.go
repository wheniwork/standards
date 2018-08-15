package data

import (
	"github.com/ecourant/standards/Site/filtering"
)

type DSessionType interface {
	Constraints() filtering.RequestConstraints
}

type DSession struct {
	UserID    int
	IsManager bool
}

func (ctx DSession) Users() DUsers {
	return DUsers{ctx}
}

func (ctx DSession) Shifts() DShifts {
	return DShifts{ctx}
}