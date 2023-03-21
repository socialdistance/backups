package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sync"
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

func NewClient(app internalapp.App, logger internalapp.Logger, workerUuid uuid.UUID) *Client {
	return &Client{
		logger:     logger,
		app:        app,
		workerUuid: workerUuid,
	}
}

func (c *Client) RequestToControlServer() (*ResponseTask, error) {
	// TODO: refactor
	taskInfo, err := internalstorage.NewTask(c.workerUuid)
	if err != nil {
		c.logger.Error("TaskInfo can't create object with err:", zap.Error(err))
		return nil, err
	}

	url := fmt.Sprintf("http://localhost:8080/api/command?id=%s&address=%s&command=%s&hostname=%s", taskInfo.WorkerUuid, taskInfo.Address, taskInfo.Command, taskInfo.Hostname)

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
	c.logger.Info("Client: response body:", zap.ByteString("body", resBody))

	responseTask := ResponseTask{}
	err = json.Unmarshal(resBody, &responseTask)
	if err != nil {
		c.logger.Error("Client could not parse to struct", zap.Error(err))
		return nil, err
	}

	return &responseTask, nil
}

func (c *Client) ExecuteBackupScriptClient(wg *sync.WaitGroup) error {
	defer wg.Done()
	err := c.app.ExecuteBackupScript("backup.sh")
	if err != nil {
		c.logger.Error("Error execute bash script:", zap.Error(err))
		return err
	}

	return nil
}

func (c *Client) SendFile(wg *sync.WaitGroup) error {
	fileNameBackup := fmt.Sprintf("/home/user/backup/backup-%d-%02d-%d.tar.gz", time.Now().Year(), time.Now().Month(), time.Now().Day())
	defer wg.Done()
	err := c.app.PostFile(fileNameBackup, "http://localhost:8080/api/upload")
	if err != nil {
		c.logger.Error("Error upload file:", zap.Error(err))
		return err
	}

	return nil
}

func (c *Client) SendBackupToControlServer() error {
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() {
		if err := c.ExecuteBackupScriptClient(wg); err != nil {
			return
		}
		if err := c.SendFile(wg); err != nil {
			return
		}
	}()

	wg.Wait()
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
					fmt.Println("Manual ticker")
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

	cronTicker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-cronTicker.C:
				fmt.Println("Cron ticker")
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
