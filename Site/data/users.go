package data

import (
	"github.com/jinzhu/gorm"
	"github.com/ecourant/standards/Site/filtering"
	_ "github.com/lib/pq"
	"github.com/ecourant/standards/Site/conf"
)

var (
	UserConstraints = filtering.GenerateConstraints(User{})
)

type DUsers struct {
	DSession
}

func (ctx DUsers) Constraints() filtering.RequestConstraints {
	return UserConstraints
}

type User struct {
	ID        *int    `json:"id,omitempty" query:"15" name:"ID"`
	Name      *string `json:"name,omitempty" query:"11" name:"Name"`
	Email     *string `json:"email,omitempty" query:"11" name:"Email"`
	Phone     *string `json:"phone,omitempty" query:"11" name:"Phone"`
	CreatedAt *string `json:"created_at,omitempty" query:"11" name:"Created At"`
	UpdatedAt *string `json:"updated_at,omitempty" query:"11" name:"Updated At"`
}

func (ctx DUsers) GetUsers(params filtering.RequestParams) ([]User, *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	if err != nil {
		return nil, NewServerError("Error, could not retrieve users at this time.", err)
	}
	defer db.Close()

	result := make([]User, 0)

	db = db.
		Table("public.vw_users_api").
		Select(params.Fields).
		Order(params.Sorts).
		Offset((params.Page * params.PageSize) - params.PageSize).
		Limit(params.PageSize)

	if len(params.Filters) > 0 {
		db = filtering.WhereFilters(db, params, ctx.Constraints())
	}

	db.Scan(&result)
	return result, nil
}
