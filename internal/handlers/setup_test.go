package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getTestRoutes(r ContactRepository) http.Handler {
	mux := chi.NewRouter()

	mux.Route("/contacts", func(router chi.Router) {
		router.Get("/", ListContacts(r))
		router.Get("/search", SearchContactByName(r))
		router.Post("/", AddContact(r))
		router.Post("/update", UpdateContact(r))
		router.Delete("/", DeleteContact(r))
	})

	return mux
}
