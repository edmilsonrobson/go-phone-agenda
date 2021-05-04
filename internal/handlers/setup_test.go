package handlers

import (
	"net/http"

	"github.com/edmilsonrobson/go-phone-agenda/internal/repositories"
	"github.com/go-chi/chi/v5"
)

func getTestRoutes() http.Handler {
	mux := chi.NewRouter()

	r := repositories.NewInMemoryContactRepository()

	mux.Route("/contacts", func(router chi.Router) {
		router.Get("/", ListContacts(r))
		router.Get("/search", SearchContactByName(r))
		router.Post("/", AddContact(r))
		router.Post("/update", UpdateContact(r))
		router.Delete("/", DeleteContact(r))
	})

	return mux
}
