package main

import (
	"encoding/json"
	"net/http"
)

func (app *Application) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	name, err := savePhoto(file)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(struct {
		Url string `json:"url"`
	}{Url: name}); err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
