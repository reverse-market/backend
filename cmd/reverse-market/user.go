package main

import (
	"encoding/json"
	"net/http"
)

func (app *Application) getUser(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	user, err := app.users.GetByID(r.Context(), id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) updateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	user, err := app.users.GetByID(r.Context(), id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	if err := app.users.Update(r.Context(), user); err != nil {
		app.serverError(w, err)
	}
}
