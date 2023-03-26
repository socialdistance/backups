package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

type ResponseTask struct {
	ID          uuid.UUID
	Command     string
	Worker_UUID uuid.UUID
	Timestamp   time.Time
}

const targetURL = "http://localhost:8080"

func NewClient(app internalapp.App, logger internalapp.Logger, workerUuid uuid.UUID) *Client {
	return &Client{
		logger:     logger,
		app:        app,
		workerUuid: workerUuid,
	}
}

func (c *Client) RequestToControlServer() (*ResponseTask, error) {
	taskInfo, err := internalstorage.NewTask(c.workerUuid)
	if err != nil {
		c.logger.Error("TaskInfo can't create object with err:", zap.Error(err))
		return nil, err
	}

	url := fmt.Sprintf("%s/api/command?id=%s&address=%s&command=%s&hostname=%s", targetURL, taskInfo.WorkerUuid, taskInfo.Address, taskInfo.Command, taskInfo.Hostname)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("Client: could not create request:", zap.Error(err))
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error("Client: error making http request:", zap.Error(err))
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.logger.Error("Client: could not read response body:", zap.Error(err))
		return nil, err

	}
	c.logger.Info("Client: status code:", zap.String("body", res.Status))
	c.logger.Info("Client: response body:", zap.ByteString("body", resBody))

	responseTask := ResponseTask{}
	err = json.Unmarshal(resBody, &responseTask)
	if err != nil {
		c.logger.Error("Client could not parse to struct", zap.Error(err))
		return nil, err
	}

	return &responseTask, nil
}

func (c *Client) ExecuteBackupScriptClient(doneCh chan struct{}) chan struct{} {
	go func() {
		fmt.Println("Test1")
		//err := c.app.ExecuteBackupScript("backup.sh")
		//if err != nil {
		//	c.logger.Error("Error execute bash script:", zap.Error(err))
		//	//errorCh <- err
		//}

		//doneCh <- struct{}{}
		<-doneCh
	}()

	return nil
}

func (c *Client) SendFile(doneCh chan struct{}) chan struct{} {
	go func() {
		fmt.Println("test2")
		//fileNameBackup := fmt.Sprintf("/home/user/backup/backup-%d-%02d-%d.tar.gz", time.Now().Year(), time.Now().Month(), time.Now().Day())
		//err := c.app.PostFile(fileNameBackup, fmt.Sprintf("%s/api/upload", targetURL))
		//if err != nil {
		//	c.logger.Error("Error upload file:", zap.Error(err))
		//}

		doneCh <- struct{}{}
		//<-doneCh
	}()

	return nil
}

func (c *Client) SendBackupToControlServer() error {
	doneCh := make(chan struct{})

	err := c.ExecuteBackupScriptClient(doneCh)
	if err != nil {
		fmt.Println("ERR1:", err)
	}

	err = c.SendFile(doneCh)
	if err != nil {
		fmt.Println("err2", err)
	}

	return nil
}

func (c *Client) Run(ctx context.Context) error {
	requestTicker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-requestTicker.C:
				responseTask, err := c.RequestToControlServer()
				if err != nil {
					return
				}
				switch responseTask.Command {
				case "manual":
					err = c.SendBackupToControlServer()
					if err != nil {
						return
					}
				}
			case <-ctx.Done():
				requestTicker.Stop()
			}
		}
	}()

	cronTicker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-cronTicker.C:
				err := c.SendBackupToControlServer()
				if err != nil {
					return
				}
			case <-ctx.Done():
				cronTicker.Stop()
			}
		}
	}()

	return nil
}
