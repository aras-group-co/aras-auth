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

type GroupHandler struct {
	groupUseCase *usecase.GroupUseCase
	validator    *validator.Validate
}

func NewGroupHandler(groupUseCase *usecase.GroupUseCase) *GroupHandler {
	return &GroupHandler{
		groupUseCase: groupUseCase,
		validator:    validator.New(),
	}
}

func (h *GroupHandler) RegisterRoutes(r chi.Router) {
	r.Route("/groups", func(r chi.Router) {
		r.Post("/", h.CreateGroup)
		r.Get("/", h.ListGroups)
		r.Get("/{id}", h.GetGroup)
		r.Put("/{id}", h.UpdateGroup)
		r.Delete("/{id}", h.DeleteGroup)
		r.Post("/{id}/members", h.AddMember)
		r.Delete("/{id}/members/{userId}", h.RemoveMember)
		r.Get("/{id}/members", h.GetMembers)
	})
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	group, err := h.groupUseCase.CreateGroup(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "creation_failed", err)
		return
	}

	WriteSuccess(w, group, "Group created successfully")
}

func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
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

	response, err := h.groupUseCase.ListGroups(r.Context(), page, limit)
	if err != nil {
		WriteInternalError(w, err)
		return
	}

	WriteSuccess(w, response, "Groups retrieved successfully")
}

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid group ID")
		return
	}

	group, err := h.groupUseCase.GetGroup(r.Context(), groupID)
	if err != nil {
		WriteNotFound(w, "Group not found")
		return
	}

	WriteSuccess(w, group, "Group retrieved successfully")
}

func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid group ID")
		return
	}

	var req domain.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	group, err := h.groupUseCase.UpdateGroup(r.Context(), groupID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "update_failed", err)
		return
	}

	WriteSuccess(w, group, "Group updated successfully")
}

func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid group ID")
		return
	}

	if err := h.groupUseCase.DeleteGroup(r.Context(), groupID); err != nil {
		WriteError(w, http.StatusBadRequest, "delete_failed", err)
		return
	}

	WriteSuccess(w, nil, "Group deleted successfully")
}

func (h *GroupHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid group ID")
		return
	}

	var req domain.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	if err := h.groupUseCase.AddMember(r.Context(), groupID, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "add_member_failed", err)
		return
	}

	WriteSuccess(w, nil, "Member added successfully")
}

func (h *GroupHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "id")
	userIDStr := chi.URLParam(r, "userId")

	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid group ID")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid user ID")
		return
	}

	if err := h.groupUseCase.RemoveMember(r.Context(), groupID, userID); err != nil {
		WriteError(w, http.StatusBadRequest, "remove_member_failed", err)
		return
	}

	WriteSuccess(w, nil, "Member removed successfully")
}

func (h *GroupHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid group ID")
		return
	}

	members, err := h.groupUseCase.GetMembers(r.Context(), groupID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "get_members_failed", err)
		return
	}

	WriteSuccess(w, members, "Members retrieved successfully")
}

