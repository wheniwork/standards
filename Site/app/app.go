package app

import (
	"github.com/kataras/iris"
	"github.com/ecourant/standards/Site/controllers"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var (
	Endpoints = []func(p iris.Party){
		controllers.Users,
		controllers.Shifts,
		controllers.Summaries,
	}
)

func App() *iris.Application {
	app := iris.Default()
	app.PartyFunc("/api", func(p iris.Party) {
		p.Use(APIMiddleware)
		// Map the endppints from the endpoints array.
		for i, endpoint := range getEndpointUrls() {
			p.PartyFunc(endpoint, Endpoints[i])
			fmt.Printf("Mapped API Endpoint: /api%s\n", endpoint)
		}
	})
	return app
}

func getEndpointUrls() []string {
	points := make([]string, len(Endpoints))
	for i, end := range Endpoints {
		// this will get the name of the function, but it will be preceded with the path in the repository
		// so we can split by "." and get the last item and that is the function name
		paths := strings.Split(runtime.FuncForPC(reflect.ValueOf(end).Pointer()).Name(), ".")
		points[i] = fmt.Sprintf("/%s", strings.ToLower(paths[len(paths)-1]))
	}
	return points
}
