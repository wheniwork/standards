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

		session := data.DSession{}

		if role, err := session.Users().GetUserRole(current_user_id); err != nil {
			ctx.StatusCode(500)
			ctx.JSON(controllers.ErrorAPIResponse{
				Message: fmt.Sprintf("Error, could not verify user ID %d", current_user_id),
			})
			return
		} else if role == nil {
			ctx.StatusCode(400)
			ctx.JSON(controllers.ErrorAPIResponse{
				Message: fmt.Sprintf("Error, could not find user with ID %d", current_user_id),
			})
			return
		} else if ctx.Method() != "GET" && *role != "manager" {
			ctx.StatusCode(403)
			ctx.JSON(controllers.ErrorAPIResponse{
				Message: "Error, as an employee you do not have permissions to make this request.",
			})
			return
		} else {
			ctx.Values().Set("Session", data.DSession{UserID: current_user_id, IsManager: *role == "manager"})
			ctx.Next()
		}
	}
}
