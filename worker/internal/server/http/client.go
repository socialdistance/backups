package http

import (
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"time"
	internalapp "worker/internal/app"
	internalstorage "worker/internal/storage"
)

type Client struct {
	logger internalapp.Logger
	app    internalapp.App

	workerUuid uuid.UUID
}

func NewClient(app internalapp.App, logger internalapp.Logger, workerUuid uuid.UUID) *Client {
	return &Client{
		logger:     logger,
		app:        app,
		workerUuid: workerUuid,
	}
}

func (c *Client) RequestToControlServer() error {
	// TODO: refactor
	taskInfo, err := internalstorage.NewTask(c.workerUuid)
	if err != nil {
		// TODO: error handling
		fmt.Println("Something wrong with taskInfo")
	}

	url := fmt.Sprintf("http://localhost:8080/api/command?id=%s&address=%s&command=%s&hostname=%s", taskInfo.WorkerUuid, taskInfo.Address, taskInfo.Command, taskInfo.Hostname)

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

func (c *Client) SendBackupToControlServer() error {
	err := c.app.ExecuteBackupScript("backup.sh")
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

func (c *Client) Run(doneCh chan struct{}) error {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := c.RequestToControlServer()
				if err != nil {
					return
				}
			case <-doneCh:
				ticker.Stop()
			}
		}
	}()

	ticker2 := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-ticker2.C:
				err := c.SendBackupToControlServer()
				if err != nil {
					return
				}
			}
		}
	}()

	return nil
}
