package tests

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/ecourant/standards/Site/app"
	"github.com/kataras/iris"
	"github.com/ecourant/standards/Site/conf"
)

const testport = 8080
var (
	baseurl = fmt.Sprintf("http://localhost:%d/api", testport)
)

func StartServer() {
	if c, err := conf.LoadConfig("test_config.json"); err != nil {
		panic(err)
	} else {
		conf.Cfg = *c
	}
	app := app.App()
	go app.Run(iris.Addr(fmt.Sprintf(":%d", testport)))
}

func GetURL(url string) (*string, *int, error){
	fullurl := fmt.Sprintf("%s/%s", baseurl, url)
	req, err := http.NewRequest("GET", fullurl, nil)
	if err != nil {
		return nil, nil, err
	} else {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		b := string(body)
		fmt.Printf("[%s] Body: %s\n", fullurl, b)
		return &b, &resp.StatusCode, nil
	}
}
