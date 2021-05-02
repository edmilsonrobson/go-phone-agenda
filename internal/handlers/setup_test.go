package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getTestRoutes() http.Handler {
	mux := chi.NewRouter()

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
