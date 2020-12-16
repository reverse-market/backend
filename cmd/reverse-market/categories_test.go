package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCategoriesEndpoint(t *testing.T) {
	app := newTestApplication()
	router := app.route(false)

	tests := []struct {
		name     string
		urlPath  string
		inBody   []byte
		wantCode int
		wantBody []byte
	}{
		{"Get all categories", "/categories", nil, http.StatusOK, []byte(`[{"id":1,"name":"Mock category","photo":"/image/smth.jpg"}]`)},
		{"Get category with wrong formatted id", "/categories/abc", nil, http.StatusBadRequest, nil},
		{"Get category with not existing id", "/categories/2", nil, http.StatusNotFound, nil},
		{"Get category", "/categories/1", nil, http.StatusOK, []byte(`{"id":1,"name":"Mock category","photo":"/image/smth.jpg"}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.urlPath, bytes.NewReader(tt.inBody))
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
