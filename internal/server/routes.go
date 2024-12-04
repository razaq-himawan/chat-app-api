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
	"github.com/razaq-himawan/chat-app-api/internal/auth"
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

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	serverRepository := repository.NewServerRepository(db)
	serverService := service.NewServerService(serverRepository)
	serverHandler := handler.NewServerHandler(serverService)

	r.Get("/health", s.healthHandler)

	r.Get("/ws", handler.HandleWebSocket)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", userHandler.HandleRegister)
		r.Post("/login", userHandler.HandleLogin)

		r.Group(func(r chi.Router) {
			r.Use(auth.AuthJWT(userService))

			r.Route("/user/{userID}", func(r chi.Router) {
				r.Get("/", userHandler.HandleGetOneUser)
				r.Put("/", userHandler.HandleUpdateUserProfile)
			})

			r.Route("/server", func(r chi.Router) {
				r.Post("/create", serverHandler.CreateServer)
			})

		})
	})

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
