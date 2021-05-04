package main

import (
	"net/http"

	"github.com/edmilsonrobson/go-phone-agenda/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Routes(r handlers.ContactRepository) http.Handler {
	mux := chi.NewRouter()

	mux.Use(LogRequest)
	r.List()

	mux.Route("/contacts", func(router chi.Router) {
		router.Get("/", handlers.ListContacts(r))
		router.Get("/search", handlers.SearchContactByName(r))
		router.Post("/", handlers.AddContact(r))
		router.Post("/update", handlers.UpdateContact(r))
		router.Delete("/", handlers.DeleteContact(r))
	})

	mux.Handle("/metrics", promhttp.Handler())

	return mux
}
