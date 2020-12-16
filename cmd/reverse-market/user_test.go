package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserEndpoint(t *testing.T) {
	app := newTestApplication()
	router := app.route(false)

	tests := []struct {
		name     string
		method   string
		inBody   []byte
		wantCode int
		wantBody []byte
	}{
		{"Get your user info", http.MethodGet, nil, http.StatusOK, []byte(`{"id":1,"name":"Ivanov Ivan","email":"ivanov@mail.ru","photo":"/image/smth1.jpg","default_address_id":1}`)},
		{"Update user with empty body", http.MethodPut, nil, http.StatusBadRequest, nil},
		{"Update user with wrong formatted body", http.MethodPut, []byte("wrong"), http.StatusBadRequest, nil},
		{"Update your user info", http.MethodPut, []byte(`{"name":"name","email":"test@mail.com","photo":"/images/smth.jpg","default_address_id":1}`), http.StatusOK, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/user", bytes.NewReader(tt.inBody))
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
