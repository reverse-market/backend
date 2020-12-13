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

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", app.getAllCategories)
		r.Get("/{categoryID}", app.getCategory)
	})

	r.Route("/tags", func(r chi.Router) {
		r.Get("/", app.getAllTags)
		r.Get("/{tagID}", app.getTag)
	})

	r.With(app.auth).Route("/", func(r chi.Router) {
		r.Get("/auth/check", app.authCheck)

		r.Route("/user", func(r chi.Router) {
			r.Get("/", app.getUser)
			r.Put("/", app.updateUser)

			r.Route("/addresses", func(r chi.Router) {
				r.Get("/", app.getUserAddresses)
				r.Post("/", app.addAddress)

				r.Route("/{addressID}", func(r chi.Router) {
					r.Get("/", app.getAddress)
					r.Put("/", app.updateAddress)
					r.Delete("/", app.deleteAddress)
				})
			})
		})
	})

	imageServer := http.FileServer(http.Dir("./images"))
	r.Get("/images/*", http.StripPrefix("/images", imageServer).ServeHTTP)

	return r
}
