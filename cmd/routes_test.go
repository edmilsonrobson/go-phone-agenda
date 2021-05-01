package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	mux := routes()

	switch mux.(type) {
	case *chi.Mux:
		// We're good!
	default:
		t.Errorf("Routes is not of expected type *chi.Mux, received type %T", mux)
	}

}
