package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/razaq-himawan/chat-app-api/internal/app/model"
	"github.com/razaq-himawan/chat-app-api/utils"
)

type UserHandler struct {
	userService model.UserService
}

func NewUserHandler(userService model.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var payload model.UserLoginPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	u, err := h.userService.CheckUserCredentials(payload)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	token, err := h.userService.LoginUser(u)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// TODO Set cookies

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var payload model.UserRegisterPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.userService.CheckIfEmailOrUsernameExists(payload)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, err)
		return
	}

	createdUser, err := h.userService.RegisterUser(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdUser)
}

// TODO handleLogout

func (h *UserHandler) HandleGetOneUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	u, err := h.userService.GetUserByIDWithProfile(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, u)
}

func (h *UserHandler) HandleUpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	var payload model.UserUpdatePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	up, err := h.userService.UpdateUserProfile(userID, payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to update user profile: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, up)
}

func (h *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	var payload model.UserDeletePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	_, err := h.userService.DeleteUser(userID, payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to delete user %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user deleted",
	})
}
