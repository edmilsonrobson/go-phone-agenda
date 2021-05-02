package handlers

import (
	"net/http"

	middlewares "github.com/edmilsonrobson/go-phone-agenda/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middlewares.LogRequest)

	mux.Get("/", ListContacts)
	mux.Route("/contacts", func(router chi.Router) {
		router.Get("/", ListContacts)
		router.Get("/search", SearchContactByName)
		router.Post("/", AddContact)
		router.Post("/update", UpdateContact)
		router.Delete("/", DeleteContact)
	})
	return mux
}
