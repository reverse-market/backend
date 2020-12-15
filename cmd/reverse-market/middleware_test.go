package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/reverse-market/backend/pkg/database/models"
	"github.com/reverse-market/backend/pkg/simpletime"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	app := newTestApplication()

	tests := []struct {
		name     string
		header   string
		wantCode int
	}{
		{"No token", "", http.StatusUnauthorized},
		{"Wrong formatting", "PLAIN_TOKEN", http.StatusUnauthorized},
		{"Invalid token", "Bearer INVALID_TOKEN", http.StatusUnauthorized},
		{"Expired token", "Bearer EXPIRED_TOKEN", http.StatusUnauthorized},
		{"Token with wrong subject", "Bearer WRONG_SUBJECT_TOKEN", http.StatusUnauthorized},
		{"Valid token", "Bearer VALID_TOKEN", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				idVal := r.Context().Value(contextKeyID)
				if idVal == nil {
					t.Error("no id present")
				}

				id, ok := idVal.(int)
				if !ok {
					t.Errorf("wrong formateed id: %v", idVal)
				}

				if id != 1 {
					t.Errorf("wrong id: wanted: 1, got: %v", id)
				}
			})

			toTest := app.auth(nextHandler)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Add("Authorization", tt.header)

			rr := httptest.NewRecorder()
			toTest.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("wrong code: want %d; got %d", tt.wantCode, rr.Code)
			}
		})
	}
}

func TestAddressCtx(t *testing.T) {
	app := newTestApplication()

	expectedAddr := &models.Address{
		ID:     1,
		UserID: 1,
		Info: models.AddressInfo{
			Name:       "Mock address 1",
			Region:     "Some region",
			City:       "Some city",
			Street:     "Mock street",
			Number:     "1",
			Building:   "1",
			Appartment: "1",
			Zip:        123,
			LastName:   "Ivanov",
			FirstName:  "Ivan",
			FatherName: "Ivanovich",
		},
	}

	tests := []struct {
		name      string
		urlPath   string
		userID    int
		checkUser bool
		wantCode  int
		wantAddr  *models.Address
	}{
		{"Wrong path", "/wrong", 1, false, http.StatusNotFound, nil},
		{"No address id", "/addresses", 1, false, http.StatusNotFound, nil},
		{"Wrong formatted id", "/addresses/abc", 1, false, http.StatusBadRequest, nil},
		{"Not existing id", "/addresses/3", 1, false, http.StatusNotFound, nil},
		{"Valid id (no user check, wrong user)", "/addresses/1", 2, false, http.StatusOK, expectedAddr},
		{"Valid id (no user check, right user)", "/addresses/1", 1, false, http.StatusOK, expectedAddr},
		{"Valid id (with user check, wrong user)", "/addresses/1", 2, true, http.StatusForbidden, nil},
		{"Valid id (with user check, right user)", "/addresses/1", 1, true, http.StatusOK, expectedAddr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				addrVal := r.Context().Value(contextKeyAddress)
				if addrVal == nil {
					t.Error("no address present")
				}

				addr, ok := addrVal.(*models.Address)
				if !ok {
					t.Errorf("wrong type address: %v", addr)
				}

				if !reflect.DeepEqual(addr, tt.wantAddr) {
					t.Errorf("wrong address: want %v; got %v", tt.wantAddr, addr)
				}
			})

			toTest := chi.NewMux()
			toTest.With(app.addressCtx(tt.checkUser)).Get("/addresses/{addressID}", nextHandler)

			req := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)
			req = req.WithContext(context.WithValue(req.Context(), contextKeyID, tt.userID))

			rr := httptest.NewRecorder()
			toTest.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("wrong code: want %d; got %d", tt.wantCode, rr.Code)
			}
		})
	}
}

func TestRequestCtx(t *testing.T) {
	app := newTestApplication()

	expectedReq := &models.Request{
		ID:          1,
		UserID:      1,
		Username:    "Ivanov Ivan",
		CategoryID:  1,
		Name:        "name 1",
		ItemName:    "item name 1",
		Description: "description 1",
		Photos:      []string{},
		Price:       100,
		Quantity:    1,
		Date:        simpletime.SimpleTime(time.Date(2020, 12, 14, 0, 0, 0, 0, time.UTC)),
		Finished:    false,
		Tags: []*models.TagInRequest{
			{
				ID:   1,
				Name: "Синий",
			},
			{
				ID:   2,
				Name: "Крансый",
			},
		},
	}

	tests := []struct {
		name      string
		urlPath   string
		userID    int
		checkUser bool
		wantCode  int
		wantReq   *models.Request
	}{
		{"Wrong path", "/wrong", 1, false, http.StatusNotFound, nil},
		{"No request id", "/requests", 1, false, http.StatusNotFound, nil},
		{"Wrong formatted id", "/requests/abc", 1, false, http.StatusBadRequest, nil},
		{"Not existing id", "/requests/4", 1, false, http.StatusNotFound, nil},
		{"Valid id (no user check, wrong user)", "/requests/1", 2, false, http.StatusOK, expectedReq},
		{"Valid id (no user check, right user)", "/requests/1", 1, false, http.StatusOK, expectedReq},
		{"Valid id (with user check, wrong user)", "/requests/1", 2, true, http.StatusForbidden, nil},
		{"Valid id (with user check, right user)", "/requests/1", 1, true, http.StatusOK, expectedReq},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqVal := r.Context().Value(contextKeyRequest)
				if reqVal == nil {
					t.Error("no request present")
				}

				req, ok := reqVal.(*models.Request)
				if !ok {
					t.Errorf("wrong type request: %v", req)
				}

				if !reflect.DeepEqual(req, tt.wantReq) {
					t.Errorf("wrong request: want %v; got %v", tt.wantReq, req)
				}
			})

			toTest := chi.NewMux()
			toTest.With(app.requestCtx(tt.checkUser)).Get("/requests/{requestID}", nextHandler)

			req := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)
			req = req.WithContext(context.WithValue(req.Context(), contextKeyID, tt.userID))

			rr := httptest.NewRecorder()
			toTest.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("wrong code: want %d; got %d", tt.wantCode, rr.Code)
			}
		})
	}
}

func TestProposalCtx(t *testing.T) {
	app := newTestApplication()

	expectedProposal := &models.Proposal{
		ID:          1,
		UserID:      1,
		Username:    "Petr Petrovich",
		RequestID:   1,
		Name:        "name 1",
		ItemName:    "item name 1",
		Description: "description 1",
		Photos:      []string{},
		Price:       100,
		Quantity:    1,
		Date:        simpletime.SimpleTime(time.Date(2020, 12, 16, 0, 0, 0, 0, time.UTC)),
		BoughtById:  nil,
	}

	tests := []struct {
		name         string
		urlPath      string
		userID       int
		checkUser    bool
		wantCode     int
		wantProposal *models.Proposal
	}{
		{"Wrong path", "/wrong", 1, false, http.StatusNotFound, nil},
		{"No proposal id", "/proposals", 1, false, http.StatusNotFound, nil},
		{"Wrong formatted id", "/proposals/abc", 1, false, http.StatusBadRequest, nil},
		{"Not existing id", "/requests/4", 1, false, http.StatusNotFound, nil},
		{"Valid id (no user check, wrong user)", "/proposals/1", 2, false, http.StatusOK, expectedProposal},
		{"Valid id (no user check, right user)", "/proposals/1", 1, false, http.StatusOK, expectedProposal},
		{"Valid id (with user check, wrong user)", "/proposals/1", 2, true, http.StatusForbidden, nil},
		{"Valid id (with user check, right user)", "/proposals/1", 1, true, http.StatusOK, expectedProposal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				proposalVal := r.Context().Value(contextKeyProposal)
				if proposalVal == nil {
					t.Error("no proposal present")
				}

				proposal, ok := proposalVal.(*models.Proposal)
				if !ok {
					t.Errorf("wrong type proposal: %v", proposal)
				}

				if !reflect.DeepEqual(proposal, tt.wantProposal) {
					t.Errorf("wrong proposal: want %v; got %v", tt.wantProposal, proposal)
				}
			})

			toTest := chi.NewMux()
			toTest.With(app.proposalCtx(tt.checkUser)).Get("/proposals/{proposalID}", nextHandler)

			req := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)
			req = req.WithContext(context.WithValue(req.Context(), contextKeyID, tt.userID))

			rr := httptest.NewRecorder()
			toTest.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("wrong code: want %d; got %d", tt.wantCode, rr.Code)
			}
		})
	}
}
