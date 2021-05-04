package main

import (
	"net/http"

	"github.com/edmilsonrobson/go-phone-agenda/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(LogRequest)

	mux.Get("/", handlers.ListContacts)
	mux.Route("/contacts", func(router chi.Router) {
		router.Get("/", handlers.ListContacts)
		router.Get("/search", handlers.SearchContactByName)
		router.Post("/", handlers.AddContact)
		router.Post("/update", handlers.UpdateContact)
		router.Delete("/", handlers.DeleteContact)
	})

	mux.Handle("/metrics", promhttp.Handler())

	return mux
}
