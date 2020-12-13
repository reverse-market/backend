package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/reverse-market/backend/pkg/database/models"
	"net/http"
	"strconv"
)

func (app *Application) getAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := app.categories.GetAll(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) getCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "categoryID"))
	if err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	category, err := app.categories.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, err, http.StatusNotFound)
			return
		}
		app.serverError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(category); err != nil {
		app.serverError(w, err)
	}
}
