package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"strings"

	"github.com/reverse-market/backend/pkg/database/models"
)

type contextKey string

const (
	contextKeyID       = contextKey("id")
	contextKeyAddress  = contextKey("address")
	contextKeyRequest  = contextKey("request")
	contextKeyProposal = contextKey("proposal")
)

var (
	ErrCantRetrieveID        = errors.New("can't retrieve id")
	AnotherUserAddressError  = errors.New("another user's address")
	AnotherUserRequestError  = errors.New("another user's request")
	AnotherUserProposalError = errors.New("another user's proposal")
)

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

		if _, err := app.users.GetByID(r.Context(), id); err != nil {
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

func (app *Application) addressCtx(checkUser bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(chi.URLParam(r, "addressID"))
			if err != nil {
				app.clientError(w, err, http.StatusBadRequest)
				return
			}

			address, err := app.addresses.GetByID(r.Context(), id)
			if err != nil {
				if errors.Is(err, models.ErrNoRecord) {
					app.clientError(w, err, http.StatusNotFound)
					return
				}
				app.serverError(w, err)
				return
			}

			if checkUser {
				userID, ok := r.Context().Value(contextKeyID).(int)
				if !ok {
					app.serverError(w, ErrCantRetrieveID)
					return
				}

				if address.UserID != userID {
					app.clientError(w, AnotherUserAddressError, http.StatusForbidden)
					return
				}
			}

			addressCtx := context.WithValue(r.Context(), contextKeyAddress, address)
			next.ServeHTTP(w, r.WithContext(addressCtx))
		})
	}
}

func (app *Application) requestCtx(checkUser bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(chi.URLParam(r, "requestID"))
			if err != nil {
				app.clientError(w, err, http.StatusBadRequest)
				return
			}

			request, err := app.requests.GetByID(r.Context(), id)
			if err != nil {
				if errors.Is(err, models.ErrNoRecord) {
					app.clientError(w, err, http.StatusNotFound)
					return
				}
				app.serverError(w, err)
				return
			}

			if checkUser {
				userID, ok := r.Context().Value(contextKeyID).(int)
				if !ok {
					app.serverError(w, ErrCantRetrieveID)
					return
				}

				if request.UserID != userID {
					app.clientError(w, AnotherUserRequestError, http.StatusForbidden)
					return
				}
			}

			requestCtx := context.WithValue(r.Context(), contextKeyRequest, request)
			next.ServeHTTP(w, r.WithContext(requestCtx))
		})
	}
}

func (app *Application) proposalCtx(checkUser bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(chi.URLParam(r, "proposalID"))
			if err != nil {
				app.clientError(w, err, http.StatusBadRequest)
				return
			}

			proposal, err := app.proposals.GetByID(r.Context(), id)
			if err != nil {
				if errors.Is(err, models.ErrNoRecord) {
					app.clientError(w, err, http.StatusNotFound)
					return
				}
				app.serverError(w, err)
				return
			}

			if checkUser {
				userID, ok := r.Context().Value(contextKeyID).(int)
				if !ok {
					app.serverError(w, ErrCantRetrieveID)
					return
				}

				if proposal.UserID != userID {
					app.clientError(w, AnotherUserProposalError, http.StatusForbidden)
					return
				}
			}

			proposalCtx := context.WithValue(r.Context(), contextKeyProposal, proposal)
			next.ServeHTTP(w, r.WithContext(proposalCtx))
		})
	}
}
