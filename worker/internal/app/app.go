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
	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
	Sync() error
}

func NewApp(logg Logger) *App {
	return &App{
		logger: logg,
	}
}

func (a *App) ExecuteBackupScript(path string) error {
	a.logger.Info("[+] Executing backup script")

	// TODO:
	_, err := exec.Command("/bin/sh", path).Output()
	if err != nil {
		fmt.Printf("error %s", err)
		return err
	}
	//output := string(cmd)
	//fmt.Println("OUTPUT", output)

	return nil
}

func (a *App) PostFile(filename string, targetUrl string) error {
	// TODO:
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

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
	return nil
}
