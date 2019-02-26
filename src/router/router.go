package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"


)

// Load loads the services, middlewares, routes, handlers.
func Load() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/api", api.NewRouter())
	r.Mount("/web", web.NewRouter())

	return r
}
