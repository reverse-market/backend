package main

import (
	"encoding/json"
	"github.com/reverse-market/backend/pkg/database/models"
	"net/http"
)

type AddressFlattened struct {
	ID int `json:"id"`
	*models.AddressInfo
}

func (app *Application) addAddress(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	address := &models.Address{UserID: userID}
	if err := json.NewDecoder(r.Body).Decode(&address.Info); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	if _, err := app.addresses.Add(r.Context(), address); err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *Application) getUserAddresses(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(contextKeyID).(int)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	addresses, err := app.addresses.GetByUserID(r.Context(), userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	flattened := make([]*AddressFlattened, len(addresses))
	for i, addr := range addresses {
		flattened[i] = &AddressFlattened{
			ID:          addr.ID,
			AddressInfo: &addr.Info,
		}
	}

	if err := json.NewEncoder(w).Encode(flattened); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) getAddress(w http.ResponseWriter, r *http.Request) {
	address, ok := r.Context().Value(contextKeyAddress).(*models.Address)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	if err := json.NewEncoder(w).Encode(&AddressFlattened{
		ID:          address.ID,
		AddressInfo: &address.Info,
	}); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) updateAddress(w http.ResponseWriter, r *http.Request) {
	address, ok := r.Context().Value(contextKeyAddress).(*models.Address)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&address.Info); err != nil {
		app.clientError(w, err, http.StatusBadRequest)
		return
	}

	if err := app.addresses.Update(r.Context(), address); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) deleteAddress(w http.ResponseWriter, r *http.Request) {
	address, ok := r.Context().Value(contextKeyAddress).(*models.Address)
	if !ok {
		app.serverError(w, ErrCantRetrieveID)
		return
	}

	if err := app.addresses.Delete(r.Context(), address.ID); err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
