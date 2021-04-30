package main

import (
	"net/http"

	"github.com/edmilsonrobson/go-phone-agenda/handlers"
	"github.com/go-chi/chi/v5"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", handlers.ListContacts)
	mux.Route("/contacts", func(router chi.Router) {
		router.Get("/", handlers.ListContacts)
		router.Post("/", handlers.AddContact)
	})
	return mux
}
