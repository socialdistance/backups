package app

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
)

type App struct {
	logger Logger
}

type Logger interface {
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
}

func NewApp(logg Logger) *App {
	return &App{
		logger: logg,
	}
}

func (a *App) ExecuteBackupScript(path string) error {
	a.logger.Info("[+] Executing backup script")

	_, err := exec.Command("/bin/sh", path).Output()
	if err != nil {
		a.logger.Error("Error execute backup script:", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) PostFile(filename string, targetUrl string) error {
	// TODO:
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("files", filename)
	if err != nil {
		a.logger.Error("Error writing to buffer:", zap.Error(err))

		return err
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		a.logger.Error("Error opening file:", zap.Error(err))
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		a.logger.Error("Error copy file:", zap.Error(err))
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		a.logger.Error("Error send to targetUrl:", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		a.logger.Error("Error read body:", zap.Error(err))
		return err
	}

	fmt.Println(resp.Status)
	a.logger.Info("[+] Status code:", zap.String("status", resp.Status))
	a.logger.Info("[+] Response:", zap.String("status", string(respBody)))
	return nil
}
