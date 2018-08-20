package tests

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/ECourant/standards/app"
	"github.com/kataras/iris"
	"bytes"
)

const testport = 8080

var (
	baseurl = fmt.Sprintf("http://localhost:%d/api", testport)
)

func StartServer() {

	app := app.App()
	go app.Run(iris.Addr(fmt.Sprintf(":%d", testport)))
}

func GetURL(url string) (*string, *int, error) {
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
		fmt.Printf("[GET %s %d] Body: %s\n", fullurl, resp.StatusCode, b)
		return &b, &resp.StatusCode, nil
	}
}

func PostURL(url string, request string) (*string, *int, error) {
	fullurl := fmt.Sprintf("%s/%s", baseurl, url)
	req, err := http.NewRequest("POST", fullurl, bytes.NewBuffer([]byte(request)))
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
		fmt.Printf("[POST %s %d] Body: %s\n", fullurl, resp.StatusCode, b)
		return &b, &resp.StatusCode, nil
	}
}

func PutURL(url string, request string) (*string, *int, error) {
	fullurl := fmt.Sprintf("%s/%s", baseurl, url)
	req, err := http.NewRequest("PUT", fullurl, bytes.NewBuffer([]byte(request)))
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
		fmt.Printf("[PUT %s %d] Body: %s\n", fullurl, resp.StatusCode, b)
		return &b, &resp.StatusCode, nil
	}
}
