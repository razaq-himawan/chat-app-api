package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/razaq-himawan/chat-app-api/internal/app/model"
	"github.com/razaq-himawan/chat-app-api/internal/auth"
	"github.com/razaq-himawan/chat-app-api/utils"
)

type ServerHandler struct {
	serverService model.ServerService
	userService   model.UserService
}

func NewServerHandler(
	serverService model.ServerService,
	userService model.UserService,
) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
		userService:   userService,
	}
}

func (h *ServerHandler) RegisterRoutes(r chi.Router) {
	r.HandleFunc("POST /create-server", auth.AuthJWT(h.CreateServer, h.userService))
}

func (h *ServerHandler) CreateServer(w http.ResponseWriter, r *http.Request) {
	var payload model.CreateServerPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	ctx := r.Context()
	userID := auth.GetUserIDFromContext(ctx)
	log.Printf("UserId: %v", userID)

	createdServer, err := h.serverService.CreateServerWithMembersAndChannels(payload, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdServer)
}
