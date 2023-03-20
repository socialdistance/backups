package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
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
		c.logger.Error("TaskInfo can't create object with err:", zap.Error(err))
		return err
	}

	url := fmt.Sprintf("http://localhost:8080/api/command?id=%s&address=%s&command=%s&hostname=%s", taskInfo.WorkerUuid, taskInfo.Address, taskInfo.Command, taskInfo.Hostname)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("Client: could not create request:", zap.Error(err))
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error("Client: error making http request:", zap.Error(err))
		return err
	}

	c.logger.Info("client: got response!")
	c.logger.Info("client: status code", zap.Int("status code", res.StatusCode))

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.logger.Error("Client: could not read response body:", zap.Error(err))
		return err
	}
	c.logger.Info("Client: response body:", zap.ByteString("body", resBody))

	return nil
}

func (c *Client) ExecuteBackupScriptClient(wg *sync.WaitGroup) error {
	defer wg.Done()
	err := c.app.ExecuteBackupScript("backup.sh")
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

func (c *Client) SendFile(wg *sync.WaitGroup) error {
	// TODO: filename backup
	fileNameBackup := fmt.Sprintf("/home/user/backup/backup-%d-%02d-%d.tar.gz", time.Now().Year(), time.Now().Month(), time.Now().Day())
	defer wg.Done()
	err := postFile(fileNameBackup, "http://localhost:8080/api/upload")
	if err != nil {
		fmt.Println("Error upload file:", err)
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
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := c.RequestToControlServer()
				if err != nil {
					return
				}
			case <-ctx.Done():
				ticker.Stop()
			}
		}
	}()

	ticker2 := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker2.C:
				err := c.SendBackupToControlServer()
				if err != nil {
					return
				}
			case <-ctx.Done():
				ticker.Stop()
			}
		}
	}()

	return nil
}

// TODO: replace it to app layer
func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("files", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}
