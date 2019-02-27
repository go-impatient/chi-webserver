package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/moocss/chi-webserver/src/handler/api/users"
	"github.com/moocss/chi-webserver/src/service"
)

func NewRouter(s service.Service ) http.Handler {
	r := chi.NewRouter()

	// request time out
	r.Use(middleware.Timeout(time.Second * 3))



	r.Route("/user", func(r chi.Router) {
		r.Get("/{id}", users.HandleFindById(s))
	})
	return r
}
