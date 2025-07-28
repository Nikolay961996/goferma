package router

import (
	"github.com/go-chi/chi/v5"
)

func GofermaRouter() *chi.Mux {
	router := chi.NewRouter()

	return router
}
