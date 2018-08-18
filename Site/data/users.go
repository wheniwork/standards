package data

import (
	"github.com/jinzhu/gorm"
	"github.com/ecourant/standards/Site/filtering"
	_ "github.com/lib/pq"
	"github.com/ecourant/standards/Site/conf"
	"strings"
	"fmt"
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
	Role	  *string `json:"role,omitempty" query:"11" name:"Role"`
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

func (ctx DUsers) CreateUser(user User) (response *User, rerr *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	if err != nil {
		return nil, NewServerError("Error, could not create users at this time.", err)
	}
	defer db.Close()

	db = db.Begin()
	// I hate how this looks in golang, but basically if there is a panic or an error somewhere further down, the transaction will rollback.
	defer func() {
		if r := recover(); r != nil {
			db.Rollback()
			response = nil
			rerr = NewServerError("Error, could not create user at this time.", err)
			return
		}
	}()

	if err := ctx.Users().verifyUser(nil, user, db); err != nil {
		return nil, err
	}

	result := make([]User, 0)

	if err := db.Raw(`
		INSERT INTO public.users (name, email, phone, role)
		VALUES (?, ?, ?, ?)
		RETURNING id,
				  name,
   				  email,
				  phone,
				  role,
				  to_char(created_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS created_at,
				  to_char(updated_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS updated_at;
	`, user.Name, user.Email, user.Phone, user.Role).Scan(&result).Error; err != nil {
		db.Rollback()
		return nil, NewServerError("Error, could not create user at this time.", err)
	}
	db.Commit()
	return &result[0], nil
}

func (ctx DUsers) UpdateUser(id int, user User) (response *User, rerr *DError) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	if err != nil {
		return nil, NewServerError("Error, could not update users at this time.", err)
	}
	defer db.Close()

	db = db.Begin()
	// I hate how this looks in golang, but basically if there is a panic or an error somewhere further down, the transaction will rollback.
	defer func() {
		if r := recover(); r != nil {
			db.Rollback()
			response = nil
			rerr = NewServerError("Error, could not update user at this time.", err)
			return
		}
	}()

	if err := ctx.Users().verifyUser(&id, user, db); err != nil {
		return nil, err
	}

	result := make([]User, 0)

	if err := db.Raw(`
		UPDATE public.users SET
			name=COALESCE(?, name),
			email=COALESCE(?, email),
			phone=COALESCE(?, phone),
			role=COALESCE(?, role),
			updated_at=LOCALTIMESTAMP
		WHERE id = ?
		RETURNING id,
				  name,
   				  email,
				  phone,
				  role,
				  to_char(created_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS created_at,
				  to_char(updated_at, 'Dy, Mon DD HH24:MI:SS.MS YYYY') AS updated_at;
	`, user.Name, user.Email, user.Phone, user.Role, user.ID).Scan(&result).Error; err != nil {
		db.Rollback()
		return nil, NewServerError("Error, could not update user at this time.", err)
	}
	db.Commit()
	return &result[0], nil
}

func (ctx DUsers) GetUserRole(id int) (*string, error) {
	db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	result := struct{
		Role string
	}{}

	db = db.Table("public.vw_users_api").
		Select("role").
		Where("id = ?", id).
		First(&result)

	if result.Role == "" {
		return nil, nil
	} else {
		return &result.Role, nil
	}
}

func (ctx DUsers) verifyUser(id *int, user User, db *gorm.DB) *DError {
	roles := map[string]bool{
		"employee": true,
		"manager": true,
	}
	if id != nil {
		if user.Role != nil {
			if _, ok := roles[strings.ToLower(*user.Role)]; !ok {
				return NewClientError(fmt.Sprintf("Error, role (%s) is not valid, must be `employee` or `manager`.", *user.Role), nil)
			}
		}
		if user.Name != nil && len(strings.TrimSpace(*user.Name)) == 0 {
			return NewClientError("Error, name cannot be blank.", nil)
		}
	} else {
		if (user.Email == nil || len(strings.TrimSpace(*user.Email)) == 0) && (user.Phone == nil || len(strings.TrimSpace(*user.Phone)) == 0) {
			return NewClientError("Error, an email or a phone number is required.", nil)
		}

		if user.Name == nil || len(strings.TrimSpace(*user.Name)) == 0 {
			return NewClientError("Error, name cannot be blank.", nil)
		}
	}
	return nil
}
