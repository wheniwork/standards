package controllers

import (
	"github.com/kataras/iris"
)

func Users(p iris.Party) {
	p.Get("/", func(ctx iris.Context) {
		ctx.JSON(struct {
			success bool
		}{ true })
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

	p.Delete("/{id:long}", func(ctx iris.Context) {
		ctx.JSON(struct {
			success bool
		}{ true })
	})
}
