package controllers

import (
	"github.com/kataras/iris"
	"github.com/ecourant/standards/Site/filtering"
	"github.com/ecourant/standards/Site/data"
)

func Shifts(p iris.Party) {
	p.Get("/", func(ctx iris.Context) {
		if params, err := filtering.ParseRequestParams(ctx, data.ShiftConstraints, filtering.StandardRequest); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Message:err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Shifts().GetShifts(*params); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success:true,
				Results:result,
			})
		}
	})

	p.Get("/mine", func(ctx iris.Context) {
		if params, err := filtering.ParseRequestParams(ctx, data.ShiftConstraints, filtering.StandardRequest); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Message:err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Shifts().GetMyShifts(*params); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success:true,
				Results:result,
			})
		}
	})

	p.Post("/", func(ctx iris.Context) {
		ctx.JSON(struct {
			success bool
		}{ true })
	})

	p.Put("/", func(ctx iris.Context) {
		ctx.JSON(struct {
			success bool
		}{ true })
	})

	p.Put("/{id:long}", func(ctx iris.Context) {
		ctx.JSON(struct {
			success bool
		}{ true })
	})

	p.Delete("/{id:long}", func(ctx iris.Context) {
		ctx.JSON(struct {
			success bool
		}{ true })
	})
}

