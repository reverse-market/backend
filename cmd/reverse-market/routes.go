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
		middleware.RedirectSlashes,
	)

	r.Post("/auth/sign_in", app.signIn)

	r.With(app.auth).Route("/", func(r chi.Router) {
		r.Get("/auth/check", app.authCheck)

		r.Route("/user", func(r chi.Router) {
			r.Get("/", app.getUser)
			r.Put("/", app.updateUser)
		})
	})

	imageServer := http.FileServer(http.Dir("./images"))
	r.Get("/images/*", http.StripPrefix("/images", imageServer).ServeHTTP)

	return r
}
