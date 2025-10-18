package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/aras-services/aras-auth/internal/domain"
	"github.com/aras-services/aras-auth/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
	validator   *validator.Validate
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		validator:   validator.New(),
	}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.RefreshToken)
		r.Post("/logout", h.Logout)
		r.Post("/verify-email", h.VerifyEmail)
		r.Post("/forgot-password", h.ForgotPassword)
		r.Post("/reset-password", h.ResetPassword)
		r.Post("/change-password", h.ChangePassword)
		r.Post("/introspect", h.IntrospectToken)
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	response, err := h.authUseCase.Register(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "registration_failed", err)
		return
	}

	WriteSuccess(w, response, "User registered successfully")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	response, err := h.authUseCase.Login(r.Context(), &req)
	if err != nil {
		WriteUnauthorized(w, "Invalid credentials")
		return
	}

	WriteSuccess(w, response, "Login successful")
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	response, err := h.authUseCase.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		WriteUnauthorized(w, "Invalid refresh token")
		return
	}

	WriteSuccess(w, response, "Token refreshed successfully")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	if err := h.authUseCase.Logout(r.Context(), req.RefreshToken); err != nil {
		WriteError(w, http.StatusBadRequest, "logout_failed", err)
		return
	}

	WriteSuccess(w, nil, "Logout successful")
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id" validate:"required,uuid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		WriteValidationError(w, "Invalid user ID")
		return
	}

	if err := h.authUseCase.VerifyEmail(r.Context(), userID); err != nil {
		WriteError(w, http.StatusBadRequest, "email_verification_failed", err)
		return
	}

	WriteSuccess(w, nil, "Email verified successfully")
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req domain.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	if err := h.authUseCase.ForgotPassword(r.Context(), &req); err != nil {
		WriteError(w, http.StatusBadRequest, "forgot_password_failed", err)
		return
	}

	WriteSuccess(w, nil, "Password reset email sent")
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req domain.ConfirmResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	if err := h.authUseCase.ResetPassword(r.Context(), &req); err != nil {
		WriteError(w, http.StatusBadRequest, "password_reset_failed", err)
		return
	}

	WriteSuccess(w, nil, "Password reset successfully")
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
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

	var req domain.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	if err := h.authUseCase.ChangePassword(r.Context(), userID, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "password_change_failed", err)
		return
	}

	WriteSuccess(w, nil, "Password changed successfully")
}

func (h *AuthHandler) IntrospectToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	introspection, err := h.authUseCase.IntrospectToken(r.Context(), req.Token)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "token_introspection_failed", err)
		return
	}

	WriteSuccess(w, introspection, "Token introspection successful")
}


