package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProposalsEndpoint(t *testing.T) {
	app := newTestApplication()
	router := app.route(false)

	tests := []struct {
		name     string
		method   string
		urlPath  string
		inBody   []byte
		wantCode int
		wantBody []byte
	}{
		{"Post empty body", http.MethodPost, "/proposals", nil, http.StatusBadRequest, nil},
		{"Post wrong formatted body", http.MethodPost, "/proposals", []byte("wrong"), http.StatusBadRequest, nil},
		{"Post correct body", http.MethodPost, "/proposals", []byte(`{"request_id":2,"description":"desc","photos":["/images/ssdgdg.jpg"],"price":10000,"quantity":1,"date":"15.12.2020"}`), http.StatusCreated, nil},
		{"Get proposal with wrong formatted id", http.MethodGet, "/proposals/abc", nil, http.StatusBadRequest, nil},
		{"Get proposal with not existing id", http.MethodGet, "/proposals/0", nil, http.StatusNotFound, nil},
		{"Get proposal", http.MethodGet, "/proposals/1", nil, http.StatusOK, []byte(`{"id":1,"username":"Petr Petrovich","request_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"16.12.2020"}`)},
		{"Update proposal with empty body", http.MethodPut, "/proposals/1", nil, http.StatusBadRequest, nil},
		{"Update other user proposal", http.MethodPut, "/proposals/2", []byte(`{"request_id":2,"name":"name","item_name":"item_name","description":"desc","photos":["/images/ssdgdg.jpg"],"price":10000,"quantity":1,"date":"15.12.2020"}`), http.StatusForbidden, nil},
		{"Update your proposal", http.MethodPut, "/proposals/1", []byte(`{"request_id":2, "description":"desc","photos":["/images/ssdgdg.jpg"],"price":10000,"quantity":1,"date":"15.12.2020"}`), http.StatusOK, nil},
		{"Delete other user proposal", http.MethodDelete, "/proposals/2", nil, http.StatusForbidden, nil},
		{"Delete your proposal", http.MethodDelete, "/proposals/1", nil, http.StatusOK, nil},
		{"Get bought proposals", http.MethodGet, "/user/orders_buy", nil, http.StatusOK, []byte(`[]`)},
		{"Buy empty body", http.MethodPost, "/user/orders_buy", nil, http.StatusBadRequest, nil},
		{"Buy wrong formatted body", http.MethodPost, "/user/orders_buy", []byte("wrong"), http.StatusBadRequest, nil},
		{"Buy correct body", http.MethodPost, "/user/orders_buy", []byte(`{"id": 1}`), http.StatusOK, nil},
		{"Get sold proposals", http.MethodGet, "/user/orders_sell", nil, http.StatusOK, []byte(`[{"id":3,"username":"Petr Petrovich","request_id":2,"name":"name 3","item_name":"item name 3","description":"description 3","photos":[],"price":100,"quantity":1,"date":"16.12.2020"}]`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.urlPath, bytes.NewReader(tt.inBody))
			req.Header.Add("Authorization", "Bearer VALID_TOKEN")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if tt.wantCode != recorder.Code {
				t.Errorf("wrong code: want %v got %v", tt.wantCode, recorder.Code)
			}

			body, err := ioutil.ReadAll(recorder.Body)
			if err != nil {
				t.Errorf("error reading body: %e", err)
			}

			if tt.wantBody != nil && !bytes.Contains(body, tt.wantBody) {
				t.Errorf("wrong body: want %s got %s", tt.wantBody, body)
			}
		})
	}
}
