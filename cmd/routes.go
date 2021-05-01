package main

import (
	"net/http"

	"github.com/edmilsonrobson/go-phone-agenda/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", handlers.ListContacts)
	mux.Route("/contacts", func(router chi.Router) {
		router.Get("/", handlers.ListContacts)
		router.Get("/search", handlers.SearchContactByName)
		router.Post("/", handlers.AddContact)
		router.Put("/", handlers.UpdateContact)
		router.Delete("/", handlers.DeleteContact)
	})
	return mux
}
