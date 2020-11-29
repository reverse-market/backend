package main

import (
	"fmt"
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
