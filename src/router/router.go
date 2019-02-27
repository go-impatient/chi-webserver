package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/moocss/chi-webserver/src/handler/api"
	"github.com/moocss/chi-webserver/src/handler/web"
	"github.com/moocss/chi-webserver/src/handler/api/sd"
	"github.com/moocss/chi-webserver/src/model"
	"github.com/moocss/chi-webserver/src/service"
	"github.com/moocss/chi-webserver/src/pkg/render"
)

var corsOpts = cors.Options{
	// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
	AllowedOrigins:   []string{"*"},
	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300, // Maximum value not ignored by any of major browsers
}

// NotFound creates a gin middleware for handling page not found.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	model.SendError(w)
	return
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, "Welcome to api app.", http.StatusOK)
	}
}

// Load loads the services, middlewares, routes, handlers.
func Load(s service.Service, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	// CORS
	cors := cors.New(corsOpts)
	r.Use(cors.Handler)

	// base middleware stack
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Use(middleware.RequestID)  // 每个请求的上下文中注册一个id
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middlewares...)

	r.NotFound(NotFoundHandler)

	r.Get("/", rootHandler())

	// The health check handlers
	r.Route("/sd", func(r chi.Router) {
		r.Get("/health", sd.HealthCheck())
		r.Get("/disk", sd.DiskCheck())
		r.Get("/cpu", sd.CPUCheck())
		r.Get("/ram", sd.RAMCheck())
	})

	r.Mount("/api/v1", api.NewRouter(s))
	r.Mount("/", web.NewRouter(s))

	return r
}
