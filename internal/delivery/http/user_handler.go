package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/aras-services/aras-auth/internal/domain"
	"github.com/aras-services/aras-auth/internal/usecase"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
	validator   *validator.Validate
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		validator:   validator.New(),
	}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.ListUsers)
		r.Get("/me", h.GetCurrentUser)
		r.Get("/{id}", h.GetUser)
		r.Put("/{id}", h.UpdateUser)
		r.Delete("/{id}", h.DeleteUser)
	})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	response, err := h.userUseCase.ListUsers(r.Context(), page, limit)
	if err != nil {
		WriteInternalError(w, err)
		return
	}

	WriteSuccess(w, response, "Users retrieved successfully")
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		WriteUnauthorized(w, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteUnauthorized(w, "Invalid user ID")
		return
	}

	user, err := h.userUseCase.GetCurrentUser(r.Context(), userID)
	if err != nil {
		WriteNotFound(w, "User not found")
		return
	}

	WriteSuccess(w, user, "User retrieved successfully")
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid user ID")
		return
	}

	user, err := h.userUseCase.GetUser(r.Context(), userID)
	if err != nil {
		WriteNotFound(w, "User not found")
		return
	}

	WriteSuccess(w, user, "User retrieved successfully")
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid user ID")
		return
	}

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	user, err := h.userUseCase.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "update_failed", err)
		return
	}

	WriteSuccess(w, user, "User updated successfully")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid user ID")
		return
	}

	if err := h.userUseCase.DeleteUser(r.Context(), userID); err != nil {
		WriteError(w, http.StatusBadRequest, "delete_failed", err)
		return
	}

	WriteSuccess(w, nil, "User deleted successfully")
}

