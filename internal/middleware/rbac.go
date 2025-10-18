package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	httphandler "github.com/aras-services/aras-auth/internal/delivery/http"
	"github.com/aras-services/aras-auth/internal/domain"
)

type RBACMiddleware struct {
	permissionRepo domain.PermissionRepository
}

func NewRBACMiddleware(permissionRepo domain.PermissionRepository) *RBACMiddleware {
	return &RBACMiddleware{
		permissionRepo: permissionRepo,
	}
}

func (m *RBACMiddleware) RequirePermission(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context
			userIDStr, ok := r.Context().Value("user_id").(string)
			if !ok {
				httphandler.WriteUnauthorized(w, "User not authenticated")
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				httphandler.WriteUnauthorized(w, "Invalid user ID")
				return
			}

			// Check permission
			hasPermission, err := m.permissionRepo.CheckUserPermission(userID, resource, action)
			if err != nil {
				httphandler.WriteInternalError(w, err)
				return
			}

			if !hasPermission {
				httphandler.WriteForbidden(w, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (m *RBACMiddleware) RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context
			userIDStr, ok := r.Context().Value("user_id").(string)
			if !ok {
				httphandler.WriteUnauthorized(w, "User not authenticated")
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				httphandler.WriteUnauthorized(w, "Invalid user ID")
				return
			}

			// Check if user has any of the required permissions
			hasAnyPermission := false
			for _, permission := range permissions {
				parts := strings.Split(permission, ":")
				if len(parts) != 2 {
					continue
				}

				resource, action := parts[0], parts[1]
				hasPermission, err := m.permissionRepo.CheckUserPermission(userID, resource, action)
				if err != nil {
					httphandler.WriteInternalError(w, err)
					return
				}

				if hasPermission {
					hasAnyPermission = true
					break
				}
			}

			if !hasAnyPermission {
				httphandler.WriteForbidden(w, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (m *RBACMiddleware) RequireAllPermissions(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context
			userIDStr, ok := r.Context().Value("user_id").(string)
			if !ok {
				httphandler.WriteUnauthorized(w, "User not authenticated")
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				httphandler.WriteUnauthorized(w, "Invalid user ID")
				return
			}

			// Check if user has all required permissions
			for _, permission := range permissions {
				parts := strings.Split(permission, ":")
				if len(parts) != 2 {
					httphandler.WriteForbidden(w, "Invalid permission format")
					return
				}

				resource, action := parts[0], parts[1]
				hasPermission, err := m.permissionRepo.CheckUserPermission(userID, resource, action)
				if err != nil {
					httphandler.WriteInternalError(w, err)
					return
				}

				if !hasPermission {
					httphandler.WriteForbidden(w, "Insufficient permissions")
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
