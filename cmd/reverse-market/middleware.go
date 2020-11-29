package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/reverse-market/backend/pkg/database/models"
)

type contextKey string

const contextKeyID = contextKey("id")

var ErrCantRetrieveID = errors.New("can't retrieve id")

func (app *Application) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			app.clientError(w, errors.New("no token provided"), http.StatusUnauthorized)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		id, err := app.tokens.GetIdFromToken(token)
		if err != nil {
			app.clientError(w, err, http.StatusUnauthorized)
			return
		}

		if _, err := app.users.GetById(r.Context(), id); err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.clientError(w, err, http.StatusUnauthorized)
				return
			}
			app.serverError(w, err)
			return
		}

		idContext := context.WithValue(r.Context(), contextKeyID, id)
		next.ServeHTTP(w, r.WithContext(idContext))
	})
}
