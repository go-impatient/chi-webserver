package web

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/moocss/chi-webserver/src/service"
)

func NewRouter(s service.Service) http.Handler {
	r := chi.NewRouter()

	return r
}