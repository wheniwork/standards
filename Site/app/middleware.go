package app

import (
	"github.com/kataras/iris"
	"github.com/ecourant/standards/Site/controllers"
	"github.com/ecourant/standards/Site/data"
	"github.com/jinzhu/gorm"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ecourant/standards/Site/conf"
)

func APIMiddleware(ctx iris.Context) {
	if !ctx.URLParamExists("current_user_id") {
		ctx.StatusCode(403)
		ctx.JSON(controllers.ErrorAPIResponse{
			Message: "Error, current_user_id url param must be specified!",
		})
		return
	}

	if current_user_id, err := ctx.URLParamInt("current_user_id"); err != nil {
		ctx.StatusCode(400)
		ctx.JSON(controllers.ErrorAPIResponse{
			Message: "Error, current_user_id url param is not a valid integer!",
		})
		return
	} else {
		db, err := gorm.Open("postgres", conf.Cfg.ConnectionString)
		defer db.Close()
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(controllers.ErrorAPIResponse{
				Message: err.Error(),
			})
			return
		}

		d := struct {
			Role string // I would normally just create a string variable for this, but I was unable to get GORM to work with anything but a struct for this call.
		}{}

		db.
			Table("vw_users_api").
			Select("role").
			Where("id = ?", current_user_id).
			First(&d)

		if d.Role == "" {
			ctx.StatusCode(400)
			ctx.JSON(controllers.ErrorAPIResponse{
				Message: fmt.Sprintf("Error, could not find user with id %d", current_user_id),
			})
			return
		}

		if ctx.Method() != "GET" && d.Role != "manager" {
			ctx.StatusCode(403)
			ctx.JSON(controllers.ErrorAPIResponse{
				Message: "Error, as an employee you do not have permissions to make this request.",
			})
			return
		}

		ctx.Values().Set("Session", data.DSession{UserID: current_user_id, IsManager: d.Role == "manager"})
		ctx.Next()
	}
}
