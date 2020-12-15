package main

import (
	"encoding/json"
	"github.com/reverse-market/backend/pkg/database/models"
	"github.com/reverse-market/backend/pkg/simpletime"
	"net/http"
	"time"
)

func (app *Application) addProposal(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	proposal := &models.Proposal{UserID: userID}
	if err := json.NewDecoder(r.Body).Decode(&proposal); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}
	proposal.Date = simpletime.SimpleTime(time.Now())

	if _, err := app.proposals.Add(r.Context(), proposal); err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *Application) getUserSold(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	proposals, err := app.proposals.GetByUserIDSold(r.Context(), userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(proposals); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) getUserBought(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	proposals, err := app.proposals.GetByUserIDBought(r.Context(), userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(proposals); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) getProposal(w http.ResponseWriter, r *http.Request) {
	proposal, ok := r.Context().Value(contextKeyProposal).(*models.Proposal)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	if err := json.NewEncoder(w).Encode(proposal); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) updateProposal(w http.ResponseWriter, r *http.Request) {
	proposal, ok := r.Context().Value(contextKeyProposal).(*models.Proposal)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&proposal); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	if err := app.proposals.Update(r.Context(), proposal); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) deleteProposal(w http.ResponseWriter, r *http.Request) {
	proposal, ok := r.Context().Value(contextKeyProposal).(*models.Proposal)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	if err := app.proposals.Delete(r.Context(), proposal.ID); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) buyProposal(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	idStruct := &struct {
		ID int `json:"id"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(idStruct); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	if err := app.proposals.Buy(r.Context(), idStruct.ID, userID); err != nil {
		app.serverError(w, err)
	}
}
