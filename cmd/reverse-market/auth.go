package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/reverse-market/backend/pkg/database/models"
	"google.golang.org/api/idtoken"
)

func (app *Application) authCheck(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *Application) signIn(w http.ResponseWriter, r *http.Request) {
	info := &struct {
		IDToken string `json:"id_token"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(info); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	payload, err := idtoken.Validate(r.Context(), info.IDToken, "")
	if err != nil {
		app.clientError(w, err, http.StatusUnauthorized)
		return
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		app.serverError(w, errors.New("can't retrieve email"))
		return
	}

	url, ok := payload.Claims["picture"].(string)
	if !ok {
		app.serverError(w, errors.New("can't retrieve picture url"))
		return
	}

	photoName := ""
	if url != "" {
		response, err := http.Get(url)
		if err != nil {
			app.serverError(w, err)
			return
		}
		defer response.Body.Close()

		photoName, err = savePhoto(response.Body)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	user, err := app.users.GetByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			user = &models.User{
				Name:             payload.Claims["name"].(string),
				Email:            payload.Claims["email"].(string),
				Photo:            photoName,
				DefaultAddressID: nil,
			}
			id, err := app.users.Add(r.Context(), user)
			if err != nil && !errors.Is(err, models.ErrDuplicateEmail) {
				app.serverError(w, err)
				return
			}
			user.ID = id
		} else {
			app.serverError(w, err)
			return
		}
	}

	token, err := app.tokens.CreateToken(user.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	response := &struct {
		JwtToken string `json:"jwt_token"`
	}{
		JwtToken: token,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.serverError(w, err)
	}
}
