package main

import (
	"encoding/json"
	"errors"
	"github.com/reverse-market/backend/pkg/idtoken"
	"net/http"

	"github.com/reverse-market/backend/pkg/database/models"
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

	tokenInfo, err := app.parser.Parse(r.Context(), info.IDToken)
	if err != nil {
		if errors.Is(err, idtoken.ErrInvalidToken) {
			app.clientError(w, err, http.StatusUnauthorized)
			return
		}
		app.serverError(w, err)
		return
	}

	photoName := ""
	if tokenInfo.Picture != "" {
		response, err := http.Get(tokenInfo.Picture)
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

	user, err := app.users.GetByEmail(r.Context(), tokenInfo.Email)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			user = &models.User{
				Name:             tokenInfo.Name,
				Email:            tokenInfo.Email,
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
