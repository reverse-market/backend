package main

import (
	"encoding/json"
	"errors"
	"github.com/reverse-market/backend/pkg/database/models"
	"github.com/reverse-market/backend/pkg/simpletime"
	"net/http"
	"time"
)

type RequestWithBestProposal struct {
	*models.Request
	BestProposalID *int `json:"best_proposal_id"`
}

func (app *Application) addRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	request := &models.Request{UserID: userID}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}
	request.Date = simpletime.SimpleTime(time.Now())

	if _, err := app.requests.Add(r.Context(), request); err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *Application) getUserRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	requests, err := app.requests.GetByUserID(r.Context(), userID, r.URL.Query().Get("search"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	withProposals := make([]*RequestWithBestProposal, len(requests))
	for i, req := range requests {
		withProposals[i] = &RequestWithBestProposal{Request: req}
		bestID, err := app.proposals.GetBestForRequest(r.Context(), req.ID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				continue
			}
			app.serverError(w, err)
		}
		withProposals[i].BestProposalID = &bestID
	}

	if err := json.NewEncoder(w).Encode(withProposals); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) getPublicRequest(w http.ResponseWriter, r *http.Request) {
	request, ok := r.Context().Value(contextKeyRequest).(*models.Request)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	withProposal := &RequestWithBestProposal{Request: request}
	bestID, err := app.proposals.GetBestForRequest(r.Context(), request.ID)
	if err != nil {
		if !errors.Is(err, models.ErrNoRecord) {
			app.serverError(w, err)
		}
	} else {
		withProposal.BestProposalID = &bestID
	}

	if err := json.NewEncoder(w).Encode(withProposal); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) getUserRequest(w http.ResponseWriter, r *http.Request) {
	request, ok := r.Context().Value(contextKeyRequest).(*models.Request)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	withProposal := &RequestWithBestProposal{Request: request}
	bestID, err := app.proposals.GetBestForRequest(r.Context(), request.ID)
	if err != nil {
		if !errors.Is(err, models.ErrNoRecord) {
			app.serverError(w, err)
		}
	} else {
		withProposal.BestProposalID = &bestID
	}

	if err := json.NewEncoder(w).Encode(withProposal); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) updateRequest(w http.ResponseWriter, r *http.Request) {
	request, ok := r.Context().Value(contextKeyRequest).(*models.Request)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	withProposal := &RequestWithBestProposal{Request: request}
	if err := json.NewDecoder(r.Body).Decode(&withProposal); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	if err := app.requests.Update(r.Context(), withProposal.Request); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) deleteRequest(w http.ResponseWriter, r *http.Request) {
	request, ok := r.Context().Value(contextKeyRequest).(*models.Request)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	if err := app.requests.Delete(r.Context(), request.ID); err != nil {
		app.serverError(w, err)
	}
}
