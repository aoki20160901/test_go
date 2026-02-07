// internal/router/router.go
package router

import (
	"net/http"

	"myapi/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup() http.Handler {
	r := chi.NewRouter()

	// middleware（gin.Default() 相当）
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// health check
	r.Get("/health", handler.Health)

	// v1 group
	r.Route("/v1", func(r chi.Router) {
		r.Get("/users/{id}", handler.GetUser)
		r.Post("/users", handler.CreateUser)
	})

	return r
}
