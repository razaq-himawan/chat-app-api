package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/razaq-himawan/chat-app-api/internal/app/handler"
	"github.com/razaq-himawan/chat-app-api/internal/app/repository"
	"github.com/razaq-himawan/chat-app-api/internal/app/service"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	db := s.db.GetDB()

	r.Get("/health", s.healthHandler)

	r.Get("/ws", handler.HandleWebSocket)

	r.Route("/api/v1", func(r chi.Router) {
		userRepository := repository.NewUserRepository(db)
		userService := service.NewUserService(userRepository)
		userHandler := handler.NewUserHandler(userService)
		userHandler.RegisterRoutes(r)

		serverRepository := repository.NewServerRepository(db)
		serverService := service.NewServerService(serverRepository)
		serverHandler := handler.NewServerHandler(serverService, userService)
		serverHandler.RegisterRoutes(r)
	})

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
