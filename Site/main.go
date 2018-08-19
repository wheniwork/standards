package main

import (
	"github.com/kataras/iris"
	"github.com/ecourant/standards/Site/app"
	"github.com/ecourant/standards/Site/conf"
	"fmt"
	"os"
)

func main() {
	path := "config.json"
	if os.Getenv("TRAVIS") == "true" {
		path = "config_travis.json"
		fmt.Println("Running In Travis")
	}
	if c, err := conf.LoadConfig(path); err != nil {
		panic(err)
	} else {
		conf.Cfg = *c
	}
	app := app.App()
	app.Run(iris.Addr(fmt.Sprintf(":%d", conf.Cfg.ListenPort)))
}
