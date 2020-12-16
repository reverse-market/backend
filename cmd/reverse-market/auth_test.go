package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthSignInEndpoint(t *testing.T) {
	app := newTestApplication()
	router := app.route(false)

	tests := []struct {
		name     string
		inBody   []byte
		wantCode int
		wantBody []byte
	}{
		{"Sign up with empty body", nil, http.StatusBadRequest, nil},
		{"Sign up with wrong formatted body", []byte("wrong"), http.StatusBadRequest, nil},
		{"Sign up with invalid token", []byte(`{"id_token": "INVALID_TOKEN"}`), http.StatusUnauthorized, nil},
		{"Sign up with valid token", []byte(`{"id_token": "VALID_TOKEN"}`), http.StatusOK, []byte(`{"jwt_token":"VALID_TOKEN"}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/sign_in", bytes.NewReader(tt.inBody))
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
