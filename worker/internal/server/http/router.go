package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	internalapp "worker/internal/app"
)

type Router struct {
	logger internalapp.Logger
	app    internalapp.App
}

func NewRouter(app internalapp.App, logger internalapp.Logger) *Router {
	return &Router{
		logger: logger,
		app:    app,
	}
}

func (r *Router) RequestToControlServer() error {
	// TODO: add hostname, ip address
	url := "http://localhost:8080/api/command?id=648f16fc-fdd5-4dab-84c6-e5f8852622e3"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	return nil
}
