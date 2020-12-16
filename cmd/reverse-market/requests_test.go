package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserRequestsEndpoint(t *testing.T) {
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
		{"Post empty body", http.MethodPost, "/user/requests", nil, http.StatusBadRequest, nil},
		{"Post wrong formatted body", http.MethodPost, "/user/requests", []byte("wrong"), http.StatusBadRequest, nil},
		{"Post correct body", http.MethodPost, "/user/requests", []byte(`{"category_id":1,"name":"name","item_name":"item name","description":"desc","photos":["/img/smth.jpg"],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}]}`), http.StatusCreated, nil},
		{"Get all requests", http.MethodGet, "/user/requests", nil, http.StatusOK, []byte(`[{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1},{"id":2,"username":"Ivanov Ivan","category_id":2,"name":"name 2","item_name":"item name 2","description":"description 2","photos":[],"price":1000,"quantity":1,"date":"15.12.2020","tags":[{"id":1,"name":"Синий"},{"id":3,"name":"Зелёный"}],"best_proposal_id":3}]`)},
		{"Get all requests with search", http.MethodGet, "/user/requests?search=1", nil, http.StatusOK, []byte(`[{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1}]`)},
		{"Get request with wrong formatted id", http.MethodGet, "/user/requests/abc", nil, http.StatusBadRequest, nil},
		{"Get request with not existing id", http.MethodGet, "/user/requests/0", nil, http.StatusNotFound, nil},
		{"Get other user request", http.MethodGet, "/user/requests/3", nil, http.StatusForbidden, nil},
		{"Get your request", http.MethodGet, "/user/requests/1", nil, http.StatusOK, []byte(`{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1}`)},
		{"Update request with empty body", http.MethodPut, "/user/requests/1", nil, http.StatusBadRequest, nil},
		{"Update other user request", http.MethodPut, "/user/requests/3", []byte(`{"category_id":1,"name":"name","item_name":"item name","description":"desc","photos":["/img/smth.jpg"],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}]}`), http.StatusForbidden, nil},
		{"Update your request", http.MethodPut, "/user/requests/1", []byte(`{"category_id":1,"name":"name","item_name":"item name","description":"desc","photos":["/img/smth.jpg"],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}]}`), http.StatusOK, nil},
		{"Delete other user request", http.MethodDelete, "/user/requests/3", nil, http.StatusForbidden, nil},
		{"Delete your request", http.MethodDelete, "/user/requests/1", nil, http.StatusOK, nil},
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

func TestRequestsEndpoint(t *testing.T) {
	app := newTestApplication()
	router := app.route(false)

	tests := []struct {
		name     string
		urlPath  string
		inBody   []byte
		wantCode int
		wantBody []byte
	}{
		{"Get all requests", "/requests", nil, http.StatusOK, []byte(`[{"id":3,"username":"Petrov Petr","category_id":2,"name":"name 3","item_name":"item name 3","description":"description 3","photos":[],"price":500,"quantity":1,"date":"16.12.2020","tags":[{"id":2,"name":"Крансый"},{"id":3,"name":"Зелёный"}],"best_proposal_id":null},{"id":2,"username":"Ivanov Ivan","category_id":2,"name":"name 2","item_name":"item name 2","description":"description 2","photos":[],"price":1000,"quantity":1,"date":"15.12.2020","tags":[{"id":1,"name":"Синий"},{"id":3,"name":"Зелёный"}],"best_proposal_id":3},{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1}]`)},
		{"Get all requests with pages", "/requests?page=1&size=2", nil, http.StatusOK, []byte(`[{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1}]`)},
		{"Get all requests with pages (empty page)", "/requests?page=1&size=10", nil, http.StatusOK, []byte(`[]`)},
		{"Get all requests with category", "/requests?category=2", nil, http.StatusOK, []byte(`[{"id":3,"username":"Petrov Petr","category_id":2,"name":"name 3","item_name":"item name 3","description":"description 3","photos":[],"price":500,"quantity":1,"date":"16.12.2020","tags":[{"id":2,"name":"Крансый"},{"id":3,"name":"Зелёный"}],"best_proposal_id":null},{"id":2,"username":"Ivanov Ivan","category_id":2,"name":"name 2","item_name":"item name 2","description":"description 2","photos":[],"price":1000,"quantity":1,"date":"15.12.2020","tags":[{"id":1,"name":"Синий"},{"id":3,"name":"Зелёный"}],"best_proposal_id":3}]`)},
		{"Get all requests with one tag", "/requests?tag=1", nil, http.StatusOK, []byte(`[{"id":2,"username":"Ivanov Ivan","category_id":2,"name":"name 2","item_name":"item name 2","description":"description 2","photos":[],"price":1000,"quantity":1,"date":"15.12.2020","tags":[{"id":1,"name":"Синий"},{"id":3,"name":"Зелёный"}],"best_proposal_id":3},{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1}]`)},
		{"Get all requests with multiple tags", "/requests?tag=1&tag=2", nil, http.StatusOK, []byte(`[{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1}]`)},
		{"Get all requests with prices", "/requests?price_from=200&price_to=800", nil, http.StatusOK, []byte(`[{"id":3,"username":"Petrov Petr","category_id":2,"name":"name 3","item_name":"item name 3","description":"description 3","photos":[],"price":500,"quantity":1,"date":"16.12.2020","tags":[{"id":2,"name":"Крансый"},{"id":3,"name":"Зелёный"}],"best_proposal_id":null}]`)},
		{"Get all requests with sort (date, asc)", "/requests?sort=date_asc", nil, http.StatusOK, []byte(`[{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1},{"id":2,"username":"Ivanov Ivan","category_id":2,"name":"name 2","item_name":"item name 2","description":"description 2","photos":[],"price":1000,"quantity":1,"date":"15.12.2020","tags":[{"id":1,"name":"Синий"},{"id":3,"name":"Зелёный"}],"best_proposal_id":3},{"id":3,"username":"Petrov Petr","category_id":2,"name":"name 3","item_name":"item name 3","description":"description 3","photos":[],"price":500,"quantity":1,"date":"16.12.2020","tags":[{"id":2,"name":"Крансый"},{"id":3,"name":"Зелёный"}],"best_proposal_id":null}]`)},
		{"Get all requests with sort (date, desc)", "/requests?sort=date_desc", nil, http.StatusOK, []byte(`[{"id":3,"username":"Petrov Petr","category_id":2,"name":"name 3","item_name":"item name 3","description":"description 3","photos":[],"price":500,"quantity":1,"date":"16.12.2020","tags":[{"id":2,"name":"Крансый"},{"id":3,"name":"Зелёный"}],"best_proposal_id":null},{"id":2,"username":"Ivanov Ivan","category_id":2,"name":"name 2","item_name":"item name 2","description":"description 2","photos":[],"price":1000,"quantity":1,"date":"15.12.2020","tags":[{"id":1,"name":"Синий"},{"id":3,"name":"Зелёный"}],"best_proposal_id":3},{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1}]`)},
		{"Get all requests with sort (price, asc)", "/requests?sort=price_asc", nil, http.StatusOK, []byte(`[{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1},{"id":3,"username":"Petrov Petr","category_id":2,"name":"name 3","item_name":"item name 3","description":"description 3","photos":[],"price":500,"quantity":1,"date":"16.12.2020","tags":[{"id":2,"name":"Крансый"},{"id":3,"name":"Зелёный"}],"best_proposal_id":null},{"id":2,"username":"Ivanov Ivan","category_id":2,"name":"name 2","item_name":"item name 2","description":"description 2","photos":[],"price":1000,"quantity":1,"date":"15.12.2020","tags":[{"id":1,"name":"Синий"},{"id":3,"name":"Зелёный"}],"best_proposal_id":3}]`)},
		{"Get all requests with sort (price, desc)", "/requests?sort=price_desc", nil, http.StatusOK, []byte(`[{"id":2,"username":"Ivanov Ivan","category_id":2,"name":"name 2","item_name":"item name 2","description":"description 2","photos":[],"price":1000,"quantity":1,"date":"15.12.2020","tags":[{"id":1,"name":"Синий"},{"id":3,"name":"Зелёный"}],"best_proposal_id":3},{"id":3,"username":"Petrov Petr","category_id":2,"name":"name 3","item_name":"item name 3","description":"description 3","photos":[],"price":500,"quantity":1,"date":"16.12.2020","tags":[{"id":2,"name":"Крансый"},{"id":3,"name":"Зелёный"}],"best_proposal_id":null},{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1}]`)},
		{"Get all requests with search", "/requests?search=1", nil, http.StatusOK, []byte(`[{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1}]`)},
		{"Get request with wrong formatted id", "/requests/abc", nil, http.StatusBadRequest, nil},
		{"Get request with not existing id", "/requests/0", nil, http.StatusNotFound, nil},
		{"Get request", "/requests/1", nil, http.StatusOK, []byte(`{"id":1,"username":"Ivanov Ivan","category_id":1,"name":"name 1","item_name":"item name 1","description":"description 1","photos":[],"price":100,"quantity":1,"date":"14.12.2020","tags":[{"id":1,"name":"Синий"},{"id":2,"name":"Крансый"}],"best_proposal_id":1}`)},
		{"Get prices", "/requests/prices", nil, http.StatusOK, []byte(`{"min":100,"max":1000}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.urlPath, bytes.NewReader(tt.inBody))
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
