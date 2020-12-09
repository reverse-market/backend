package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.loggers.error.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, err error, status int) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.loggers.info.Output(2, trace)

	http.Error(w, http.StatusText(status), status)
}

func (app *Application) savePhoto(url string) (string, error) {
	if url == "" {
		return "", nil
	}

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	file, err := ioutil.TempFile("images", "*.jpg")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := file.Chmod(0777); err != nil {
		return "", err
	}

	if _, err := io.Copy(file, response.Body); err != nil {
		return "", err
	}

	return fmt.Sprintf("/%s", file.Name()), nil
}
