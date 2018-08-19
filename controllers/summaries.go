package controllers


import (
	"github.com/kataras/iris"
	"github.com/ecourant/standards/filtering"
	"github.com/ecourant/standards/data"
)

func Summaries(p iris.Party) {
	p.Get("/", func(ctx iris.Context) {
		if params, err := filtering.ParseRequestParams(ctx, data.SummaryConstraints, filtering.StandardRequest); err != nil {
			ctx.StatusCode(400)
			ctx.JSON(ErrorAPIResponse{
				Message: err.Error(),
			})
		} else if result, err := ctx.Values().Get("Session").(data.DSession).Summaries().GetSummary(*params); err != nil {
			data.ErrorResponse(ctx, err)
		} else {
			ctx.JSON(APIResponse{
				Success: true,
				Results: result,
			})
		}
	})
}