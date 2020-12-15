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

	r.Route("/requests", func(r chi.Router) {
		r.Get("/", app.getPublicRequests)
		r.Get("/prices", app.getPrices)
		r.With(app.requestCtx(false)).Get("/{requestID}", app.getPublicRequest)
	})

	r.With(app.proposalCtx(false)).Get("/proposals/{proposalID}", app.getProposal)

	r.With(app.auth).Route("/", func(r chi.Router) {
		r.Get("/auth/check", app.authCheck)
		r.Post("/images", app.uploadPhoto)

		r.Route("/proposals", func(r chi.Router) {
			r.Post("/", app.addProposal)

			r.With(app.proposalCtx(true)).Route("/{proposalID}", func(r chi.Router) {
				r.Put("/", app.updateProposal)
				r.Delete("/", app.deleteProposal)
			})
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/", app.getUser)
			r.Put("/", app.updateUser)

			r.Route("/addresses", func(r chi.Router) {
				r.Get("/", app.getUserAddresses)
				r.Post("/", app.addAddress)

				r.With(app.addressCtx(true)).Route("/{addressID}", func(r chi.Router) {
					r.Get("/", app.getAddress)
					r.Put("/", app.updateAddress)
					r.Delete("/", app.deleteAddress)
				})
			})

			r.Route("/requests", func(r chi.Router) {
				r.Get("/", app.getUserRequests)
				r.Post("/", app.addRequest)

				r.With(app.requestCtx(true)).Route("/{requestID}", func(r chi.Router) {
					r.Get("/", app.getUserRequest)
					r.Put("/", app.updateRequest)
					r.Delete("/", app.deleteRequest)
				})
			})

			r.Get("/orders_sell", app.getUserSold)

			r.Route("/orders_buy", func(r chi.Router) {
				r.Get("/", app.getUserBought)
				r.Post("/", app.buyProposal)
			})
		})
	})

	imageServer := http.FileServer(http.Dir("./images"))
	r.Get("/images/*", http.StripPrefix("/images", imageServer).ServeHTTP)

	return r
}
