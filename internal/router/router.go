// internal/router/router.go
package router

import (
	"net/http"

	"myapi/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Setup() http.Handler {
	r := chi.NewRouter()

	// middleware（gin.Default() 相当）
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// health check
	r.Get("/health", handler.Health)

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./api/openapi.yaml")
	})

	// ⭐ Swagger UI
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/openapi.yaml"),
	))

	// v1 group
	r.Route("/v1", func(r chi.Router) {
		r.Get("/users/{id}", handler.GetUser)
		r.Post("/users", handler.CreateUser)
	})

	return r
}
