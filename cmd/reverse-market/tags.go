package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/reverse-market/backend/pkg/database/models"
	"net/http"
	"strconv"
)

func (app *Application) getAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := app.tags.GetAll(r.Context(), getTagFilters(r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(tags); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) getTag(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "tagID"))
	if err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	tag, err := app.tags.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, err, http.StatusNotFound)
			return
		}
		app.serverError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(tag); err != nil {
		app.serverError(w, err)
	}
}
