package data

import (
	"github.com/ecourant/standards/Site/filtering"
)

type DSessionType interface {
	Constraints() filtering.RequestConstraints
}

type DSession struct {
	UserID    int64
	IsManager bool
}

func (ctx DSession) Users() DUsers {
	return DUsers{ctx}
}