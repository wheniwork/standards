package controllers

import (
	"github.com/kataras/iris"
	"github.com/ECourant/standards/filtering"
	"github.com/ECourant/standards/data"
)

func Users(p iris.Party) {
	p.Get("/", func(ctx iris.Context) {
		if params, err := filtering.ParseRequestParams(ctx, data.UserConstraints, filtering.StandardRequest); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Message: err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Users().GetUsers(*params); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Post("/", func(ctx iris.Context) {
		newItem := data.User{}
		if err := ctx.ReadJSON(&newItem); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Success: false,
				Message: err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Users().CreateUser(newItem); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Put("/", func(ctx iris.Context) {
		newItem := data.User{}
		if err := ctx.ReadJSON(&newItem); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Success: false,
				Message: err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Users().CreateUser(newItem); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Put("/{id:long}", func(ctx iris.Context) {
		newItem := data.User{}
		if err := ctx.ReadJSON(&newItem); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Success: false,
				Message: err.Error(),
			})
		} else if id, err := ctx.Params().GetInt("id"); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Success: false,
				Message: "Error, could not parse user id.",
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Users().UpdateUser(id, newItem); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})
}
