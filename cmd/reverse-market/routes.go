package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *Application) route() http.Handler {
	r := chi.NewMux()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.RealIP,
	)

	r.Post("/auth/sign_in", app.signIn)

	r.With(app.auth).Route("/", func(r chi.Router) {
		r.Get("/auth/check", app.authCheck)

		r.Get("/user", app.getUser)
	})

	return r
}
