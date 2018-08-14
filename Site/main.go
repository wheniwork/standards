package main

import (
	"github.com/kataras/iris"
	"github.com/ecourant/standards/Site/app"
)

func main() {
	app := app.App()
	app.Run(iris.Addr(":8080"))
}
