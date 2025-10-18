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

type AuthzHandler struct {
	authzUseCase *usecase.AuthzUseCase
	validator    *validator.Validate
}

func NewAuthzHandler(authzUseCase *usecase.AuthzUseCase) *AuthzHandler {
	return &AuthzHandler{
		authzUseCase: authzUseCase,
		validator:    validator.New(),
	}
}

func (h *AuthzHandler) RegisterRoutes(r chi.Router) {
	r.Route("/roles", func(r chi.Router) {
		r.Post("/", h.CreateRole)
		r.Get("/", h.ListRoles)
		r.Get("/{id}", h.GetRole)
		r.Put("/{id}", h.UpdateRole)
		r.Delete("/{id}", h.DeleteRole)
		r.Post("/{id}/permissions", h.AssignPermissionToRole)
		r.Delete("/{id}/permissions/{permissionId}", h.RemovePermissionFromRole)
		r.Get("/{id}/permissions", h.GetRolePermissions)
	})

	r.Route("/permissions", func(r chi.Router) {
		r.Post("/", h.CreatePermission)
		r.Get("/", h.ListPermissions)
		r.Get("/{id}", h.GetPermission)
		r.Put("/{id}", h.UpdatePermission)
		r.Delete("/{id}", h.DeletePermission)
	})

	r.Route("/users/{userId}/roles", func(r chi.Router) {
		r.Post("/", h.AssignRoleToUser)
		r.Delete("/{roleId}", h.RemoveRoleFromUser)
		r.Get("/", h.GetUserRoles)
	})

	r.Route("/groups/{groupId}/roles", func(r chi.Router) {
		r.Post("/", h.AssignRoleToGroup)
		r.Delete("/{roleId}", h.RemoveRoleFromGroup)
		r.Get("/", h.GetGroupRoles)
	})

	r.Route("/authz", func(r chi.Router) {
		r.Post("/check", h.CheckPermission)
	})
}

// Role handlers
func (h *AuthzHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	role, err := h.authzUseCase.CreateRole(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "creation_failed", err)
		return
	}

	WriteSuccess(w, role, "Role created successfully")
}

func (h *AuthzHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
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

	response, err := h.authzUseCase.ListRoles(r.Context(), page, limit)
	if err != nil {
		WriteInternalError(w, err)
		return
	}

	WriteSuccess(w, response, "Roles retrieved successfully")
}

func (h *AuthzHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := chi.URLParam(r, "id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid role ID")
		return
	}

	role, err := h.authzUseCase.GetRole(r.Context(), roleID)
	if err != nil {
		WriteNotFound(w, "Role not found")
		return
	}

	WriteSuccess(w, role, "Role retrieved successfully")
}

func (h *AuthzHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := chi.URLParam(r, "id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid role ID")
		return
	}

	var req domain.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	role, err := h.authzUseCase.UpdateRole(r.Context(), roleID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "update_failed", err)
		return
	}

	WriteSuccess(w, role, "Role updated successfully")
}

func (h *AuthzHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := chi.URLParam(r, "id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid role ID")
		return
	}

	if err := h.authzUseCase.DeleteRole(r.Context(), roleID); err != nil {
		WriteError(w, http.StatusBadRequest, "delete_failed", err)
		return
	}

	WriteSuccess(w, nil, "Role deleted successfully")
}

func (h *AuthzHandler) AssignPermissionToRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := chi.URLParam(r, "id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid role ID")
		return
	}

	var req domain.AssignPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	if err := h.authzUseCase.AssignPermissionToRole(r.Context(), roleID, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "assign_permission_failed", err)
		return
	}

	WriteSuccess(w, nil, "Permission assigned to role successfully")
}

func (h *AuthzHandler) RemovePermissionFromRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := chi.URLParam(r, "id")
	permissionIDStr := chi.URLParam(r, "permissionId")

	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid role ID")
		return
	}

	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid permission ID")
		return
	}

	if err := h.authzUseCase.RemovePermissionFromRole(r.Context(), roleID, permissionID); err != nil {
		WriteError(w, http.StatusBadRequest, "remove_permission_failed", err)
		return
	}

	WriteSuccess(w, nil, "Permission removed from role successfully")
}

func (h *AuthzHandler) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleIDStr := chi.URLParam(r, "id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid role ID")
		return
	}

	permissions, err := h.authzUseCase.GetRolePermissions(r.Context(), roleID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "get_permissions_failed", err)
		return
	}

	WriteSuccess(w, permissions, "Role permissions retrieved successfully")
}

// Permission handlers
func (h *AuthzHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req domain.CreatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	permission, err := h.authzUseCase.CreatePermission(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "creation_failed", err)
		return
	}

	WriteSuccess(w, permission, "Permission created successfully")
}

func (h *AuthzHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
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

	response, err := h.authzUseCase.ListPermissions(r.Context(), page, limit)
	if err != nil {
		WriteInternalError(w, err)
		return
	}

	WriteSuccess(w, response, "Permissions retrieved successfully")
}

func (h *AuthzHandler) GetPermission(w http.ResponseWriter, r *http.Request) {
	permissionIDStr := chi.URLParam(r, "id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid permission ID")
		return
	}

	permission, err := h.authzUseCase.GetPermission(r.Context(), permissionID)
	if err != nil {
		WriteNotFound(w, "Permission not found")
		return
	}

	WriteSuccess(w, permission, "Permission retrieved successfully")
}

func (h *AuthzHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	permissionIDStr := chi.URLParam(r, "id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid permission ID")
		return
	}

	var req domain.UpdatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	permission, err := h.authzUseCase.UpdatePermission(r.Context(), permissionID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "update_failed", err)
		return
	}

	WriteSuccess(w, permission, "Permission updated successfully")
}

func (h *AuthzHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	permissionIDStr := chi.URLParam(r, "id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid permission ID")
		return
	}

	if err := h.authzUseCase.DeletePermission(r.Context(), permissionID); err != nil {
		WriteError(w, http.StatusBadRequest, "delete_failed", err)
		return
	}

	WriteSuccess(w, nil, "Permission deleted successfully")
}

// User role handlers
func (h *AuthzHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid user ID")
		return
	}

	var req domain.AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	if err := h.authzUseCase.AssignRoleToUser(r.Context(), userID, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "assign_role_failed", err)
		return
	}

	WriteSuccess(w, nil, "Role assigned to user successfully")
}

func (h *AuthzHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	roleIDStr := chi.URLParam(r, "roleId")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid user ID")
		return
	}

	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid role ID")
		return
	}

	if err := h.authzUseCase.RemoveRoleFromUser(r.Context(), userID, roleID); err != nil {
		WriteError(w, http.StatusBadRequest, "remove_role_failed", err)
		return
	}

	WriteSuccess(w, nil, "Role removed from user successfully")
}

func (h *AuthzHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid user ID")
		return
	}

	roles, err := h.authzUseCase.GetUserRoles(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "get_roles_failed", err)
		return
	}

	WriteSuccess(w, roles, "User roles retrieved successfully")
}

// Group role handlers
func (h *AuthzHandler) AssignRoleToGroup(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "groupId")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid group ID")
		return
	}

	var req domain.AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	if err := h.authzUseCase.AssignRoleToGroup(r.Context(), groupID, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "assign_role_failed", err)
		return
	}

	WriteSuccess(w, nil, "Role assigned to group successfully")
}

func (h *AuthzHandler) RemoveRoleFromGroup(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "groupId")
	roleIDStr := chi.URLParam(r, "roleId")

	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid group ID")
		return
	}

	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid role ID")
		return
	}

	if err := h.authzUseCase.RemoveRoleFromGroup(r.Context(), groupID, roleID); err != nil {
		WriteError(w, http.StatusBadRequest, "remove_role_failed", err)
		return
	}

	WriteSuccess(w, nil, "Role removed from group successfully")
}

func (h *AuthzHandler) GetGroupRoles(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "groupId")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		WriteValidationError(w, "Invalid group ID")
		return
	}

	roles, err := h.authzUseCase.GetGroupRoles(r.Context(), groupID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "get_roles_failed", err)
		return
	}

	WriteSuccess(w, roles, "Group roles retrieved successfully")
}

// Authorization check
func (h *AuthzHandler) CheckPermission(w http.ResponseWriter, r *http.Request) {
	var req domain.CheckPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteValidationError(w, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		WriteValidationError(w, err.Error())
		return
	}

	response, err := h.authzUseCase.CheckPermission(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "check_permission_failed", err)
		return
	}

	WriteSuccess(w, response, "Permission check completed")
}


