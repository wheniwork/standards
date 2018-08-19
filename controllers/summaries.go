package controllers


import (
	"github.com/kataras/iris"
	"github.com/ECourant/standards/filtering"
	"github.com/ECourant/standards/data"
)

func Summaries(p iris.Party) {
	p.Get("/", func(ctx iris.Context) {
		if params, err := filtering.ParseRequestParams(ctx, data.SummaryConstraints, filtering.StandardRequest); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Message: err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Summaries().GetSummary(nil, *params); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})

	p.Get("/{id:long}", func(ctx iris.Context) {
		if params, err := filtering.ParseRequestParams(ctx, data.SummaryConstraints, filtering.StandardRequest); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Message: err.Error(),
			})
		} else if id, err := ctx.Params().GetInt("id"); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Success:false,
				Message:"Error, could not parse user id.",
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Summaries().GetSummary(&id, *params); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})
}