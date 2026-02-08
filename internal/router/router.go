package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"myapi/internal/handler"
	"myapi/internal/service"
)

func Setup(userService *service.UserService) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handler.Health)

	r.Route("/v1", func(r chi.Router) {
		h := handler.NewUserHandler(userService)

		r.Get("/users/{id}", h.GetUser)
		r.Post("/users", h.CreateUser)
	})

	return r
}
