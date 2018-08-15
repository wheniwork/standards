package main

import (
	"github.com/kataras/iris"
	"github.com/ecourant/standards/Site/app"
	"github.com/ecourant/standards/Site/conf"
	"fmt"
)



func main() {
	if c, err := conf.LoadConfig("config.json"); err != nil { panic(err) } else {
		conf.Cfg = *c
	}
	app := app.App()
	app.Run(iris.Addr(fmt.Sprintf(":%d", conf.Cfg.ListenPort)))
}
