package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserAddressesEndpoint(t *testing.T) {
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
		{"Post empty body", http.MethodPost, "/user/addresses", nil, http.StatusBadRequest, nil},
		{"Post wrong formatted body", http.MethodPost, "/user/addresses", []byte("wrong"), http.StatusBadRequest, nil},
		{"Post correct body", http.MethodPost, "/user/addresses", []byte(`{"name":"name","region":"region","city":"city","street":"street","number":"1","building":"1","appartment":"1","zip":123,"last_name":"Ivanov","first_name":"Ivan","father_name":"Ivanovich"}`), http.StatusCreated, nil},
		{"Get all addresses", http.MethodGet, "/user/addresses", nil, http.StatusOK, []byte(`[{"id":1,"name":"Mock address 1","region":"Some region","city":"Some city","street":"Mock street","number":"1","building":"1","appartment":"1","zip":123,"last_name":"Ivanov","first_name":"Ivan","father_name":"Ivanovich"}]`)},
		{"Get other user address", http.MethodGet, "/user/addresses/2", nil, http.StatusForbidden, nil},
		{"Get your address", http.MethodGet, "/user/addresses/1", nil, http.StatusOK, []byte(`{"id":1,"name":"Mock address 1","region":"Some region","city":"Some city","street":"Mock street","number":"1","building":"1","appartment":"1","zip":123,"last_name":"Ivanov","first_name":"Ivan","father_name":"Ivanovich"}`)},
		{"Update address with empty body", http.MethodPut, "/user/addresses/1", nil, http.StatusBadRequest, nil},
		{"Update other user address", http.MethodPut, "/user/addresses/2", []byte(`{"name":"name","region":"region","city":"city","street":"street","number":"1","building":"1","appartment":"1","zip":123,"last_name":"Ivanov","first_name":"Ivan","father_name":"Ivanovich"}`), http.StatusForbidden, nil},
		{"Update your address", http.MethodPut, "/user/addresses/1", []byte(`{"name":"name","region":"region","city":"city","street":"street","number":"1","building":"1","appartment":"1","zip":123,"last_name":"Ivanov","first_name":"Ivan","father_name":"Ivanovich"}`), http.StatusOK, nil},
		{"Delete other user address", http.MethodDelete, "/user/addresses/2", nil, http.StatusForbidden, nil},
		{"Delete your address", http.MethodDelete, "/user/addresses/1", nil, http.StatusNoContent, nil},
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
