package controllers

import (
	"github.com/kataras/iris"
	"github.com/ECourant/standards/filtering"
	"github.com/ECourant/standards/data"
)

func Shifts(p iris.Party) {
	p.Get("/", func(ctx iris.Context) {
		// This is the largest repeated code that I use I think.
		// It parses the url parameters to help build the SQL query later on.
		if params, err := filtering.ParseRequestParams(ctx, data.ShiftConstraints, filtering.StandardRequest); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Message: err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Shifts().GetShifts(*params); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Get("/mine", func(ctx iris.Context) {
		if params, err := filtering.ParseRequestParams(ctx, data.ShiftConstraints, filtering.StandardRequest); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Message: err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Shifts().GetMyShifts(*params); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Get("/overlapping/{id:int}", func(ctx iris.Context) {
		if params, err := filtering.ParseRequestParams(ctx, data.ShiftConstraints, filtering.StandardRequest); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Message: err.Error(),
			})
		} else if id, err := ctx.Params().GetInt("id"); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Success: false,
				Message: "Error, could not parse shift id.",
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Shifts().GetShiftDetails(*params, id); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Get("/nonoverlapping/{id:int}/users", func(ctx iris.Context){
		if id, err := ctx.Params().GetInt("id"); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Success: false,
				Message: "Error, could not parse shift id.",
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Shifts().GetNonConflictingUsers(id); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Post("/", func(ctx iris.Context) {
		newItem := data.Shift{}
		if err := ctx.ReadJSON(&newItem); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Success: false,
				Message: err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Shifts().CreateShift(newItem); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Put("/", func(ctx iris.Context) {
		newItem := data.Shift{}
		if err := ctx.ReadJSON(&newItem); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Success: false,
				Message: err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Shifts().CreateShift(newItem); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Put("/{id:long}", func(ctx iris.Context) {
		newItem := data.Shift{}
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
				Message: "Error, could not parse shift id.",
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Shifts().UpdateShift(id, newItem); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Delete("/{id:long}", func(ctx iris.Context) {
		if id, err := ctx.Params().GetInt("id"); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Success: false,
				Message: "Error, could not parse shift id.",
			})
		} else if err := ctx.Values().Get("Session").(data.DSession).Shifts().DeleteShift(id); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
			})
		}
	})
}
