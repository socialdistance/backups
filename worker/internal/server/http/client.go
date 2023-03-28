package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	internalapp "worker/internal/app"
	internalstorage "worker/internal/storage"
)

type Client struct {
	logger    internalapp.Logger
	app       internalapp.App
	configURL string

	workerUuid uuid.UUID
}

type ResponseTask struct {
	ID         uuid.UUID
	Command    string
	WorkerUuid uuid.UUID
	Timestamp  time.Time
}

func NewClient(app internalapp.App, logger internalapp.Logger, configURL string, workerUuid uuid.UUID) *Client {
	return &Client{
		logger:     logger,
		app:        app,
		configURL:  configURL,
		workerUuid: workerUuid,
	}
}

func (c *Client) RequestToControlServer() (*ResponseTask, error) {
	taskInfo, err := internalstorage.NewTask(c.workerUuid)
	if err != nil {
		c.logger.Error("TaskInfo can't create object with err:", zap.Error(err))
		return nil, err
	}

	url := fmt.Sprintf("%s/api/command?id=%s&address=%s&command=%s&hostname=%s", c.configURL, taskInfo.WorkerUuid, taskInfo.Address, taskInfo.Command, taskInfo.Hostname)

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

func (c *Client) PostFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("files", filename)
	if err != nil {
		c.logger.Error("Error writing to buffer:", zap.Error(err))
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		c.logger.Error("Error opening file:", zap.Error(err))
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		c.logger.Error("Error copy file:", zap.Error(err))
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		c.logger.Error("Error send to targetUrl:", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Error read body:", zap.Error(err))
		return err
	}

	c.logger.Info("[+] Status code:", zap.String("status", resp.Status))
	c.logger.Info("[+] Response:", zap.String("status", string(respBody)))
	return nil
}

func (c *Client) ExecuteBackupScriptClient() error {
	err := c.app.ExecuteBackupScript("backup.sh")
	if err != nil {
		c.logger.Error("Error execute bash script:", zap.Error(err))
		return err
	}

	return nil
}

func (c *Client) SendFile() error {
	fileNameBackup := fmt.Sprintf("/home/user/backup/backup-%d-%02d-%d.tar.gz", time.Now().Year(), time.Now().Month(), time.Now().Day())
	err := c.PostFile(fileNameBackup, fmt.Sprintf("%s/api/upload", c.configURL))
	if err != nil {
		c.logger.Error("Error upload file:", zap.Error(err))
		return err
	}

	return nil
}

func (c *Client) SendBackupToControlServer() {
	err := c.ExecuteBackupScriptClient()
	if err != nil {
		c.logger.Error("Error execute bash script", zap.Error(err))
	}

	err = c.SendFile()
	if err != nil {
		c.logger.Error("Error send file", zap.Error(err))
	}
}

func (c *Client) Run(ctx context.Context) error {
	requestTicker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-requestTicker.C:
				responseTask, err := c.RequestToControlServer()
				if err != nil {
					c.logger.Error("Response task to control server have error:", zap.Error(err))
				}
				switch responseTask.Command {
				case "manual":
					c.SendBackupToControlServer()
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
				c.SendBackupToControlServer()
			case <-ctx.Done():
				cronTicker.Stop()
			}
		}
	}()

	return nil
}
