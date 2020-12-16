package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTagsEndpoint(t *testing.T) {
	app := newTestApplication()
	router := app.route(false)

	tests := []struct {
		name     string
		urlPath  string
		inBody   []byte
		wantCode int
		wantBody []byte
	}{
		{"Get all tags", "/tags", nil, http.StatusOK, []byte(`[{"id":1,"name":"Синий"},{"id":2,"name":"Красный"},{"id":3,"name":"Зелёный"}]`)},
		{"Get all tags with search", "/tags?search=синий", nil, http.StatusOK, []byte(`[{"id":1,"name":"Синий"}]`)},
		{"Get tag with wrong formatted id", "/tags/abc", nil, http.StatusBadRequest, nil},
		{"Get tag with not existing id", "/tags/0", nil, http.StatusNotFound, nil},
		{"Get tags", "/tags/1", nil, http.StatusOK, []byte(`{"id":1,"name":"Синий"}`)},
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
