package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aras-services/aras-auth/pkg/client/go/arasauth"
	"github.com/gorilla/mux"
)

// APIGateway represents an API Gateway that forwards requests to Auth service
type APIGateway struct {
	authClient *arasauth.Client
}

// NewAPIGateway creates a new API Gateway instance
func NewAPIGateway() *APIGateway {
	// Get Auth service URL from environment variable
	authURL := os.Getenv("ARAS_AUTH_URL")
	if authURL == "" {
		authURL = "http://localhost:7600" // fallback for local development
		log.Printf("ARAS_AUTH_URL not set, using default: %s", authURL)
	}

	log.Printf("Initializing API Gateway with ArasAuth URL: %s", authURL)

	return &APIGateway{
		authClient: arasauth.NewClient(authURL),
	}
}

// extractTokenFromHeader extracts Bearer token from Authorization header
func extractTokenFromHeader(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}

// HandleLogin handles user login requests from frontend
func (gw *APIGateway) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode login request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create context with timeout for auth service call
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Forward to auth service
	authResp, err := gw.authClient.Login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("Authentication failed for user %s: %v", req.Email, err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Return response to frontend
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authResp); err != nil {
		log.Printf("Failed to encode login response: %v", err)
	}
}

// HandleRegister handles user registration requests from frontend
func (gw *APIGateway) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Forward to auth service
	user, err := gw.authClient.Register(r.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// HandleChangePassword handles password change requests from frontend
func (gw *APIGateway) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	token := extractTokenFromHeader(r)
	if token == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}
	gw.authClient.SetToken(token)

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Forward to auth service
	err := gw.authClient.ChangePassword(r.Context(), req.CurrentPassword, req.NewPassword)
	if err != nil {
		http.Error(w, "Password change failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}

// HandleGetCurrentUser handles get current user requests
func (gw *APIGateway) HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	token := extractTokenFromHeader(r)
	if token == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}
	gw.authClient.SetToken(token)

	user, err := gw.authClient.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// HandleCreateRole handles role creation requests from admin panel
func (gw *APIGateway) HandleCreateRole(w http.ResponseWriter, r *http.Request) {
	token := extractTokenFromHeader(r)
	if token == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}
	gw.authClient.SetToken(token)

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Forward to auth service
	role, err := gw.authClient.CreateRole(r.Context(), req.Name, req.Description)
	if err != nil {
		http.Error(w, "Role creation failed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

// HandleListRoles handles role listing requests from admin panel
func (gw *APIGateway) HandleListRoles(w http.ResponseWriter, r *http.Request) {
	token := extractTokenFromHeader(r)
	if token == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}
	gw.authClient.SetToken(token)

	// Parse pagination parameters
	page := 1
	limit := 20
	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	roles, err := gw.authClient.ListRoles(r.Context(), page, limit)
	if err != nil {
		http.Error(w, "Failed to list roles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// HandleAssignRoleToUser handles role assignment requests
func (gw *APIGateway) HandleAssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	token := extractTokenFromHeader(r)
	if token == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}
	gw.authClient.SetToken(token)

	vars := mux.Vars(r)
	userID := vars["userId"]

	var req struct {
		RoleID string `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := gw.authClient.AssignRoleToUser(r.Context(), userID, req.RoleID)
	if err != nil {
		http.Error(w, "Role assignment failed", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Role assigned successfully"})
}

// HandleCheckPermission handles permission check requests
func (gw *APIGateway) HandleCheckPermission(w http.ResponseWriter, r *http.Request) {
	token := extractTokenFromHeader(r)
	if token == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}
	gw.authClient.SetToken(token)

	var req struct {
		UserID   string `json:"user_id"`
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hasPermission, err := gw.authClient.CheckPermission(r.Context(), req.UserID, req.Resource, req.Action)
	if err != nil {
		http.Error(w, "Permission check failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"has_permission": hasPermission})
}

// HandleTokenIntrospect handles token introspection requests
func (gw *APIGateway) HandleTokenIntrospect(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	introspection, err := gw.authClient.IntrospectToken(r.Context(), req.Token)
	if err != nil {
		http.Error(w, "Token introspection failed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(introspection)
}

// HandleHealthCheck handles health check requests
func (gw *APIGateway) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check ArasAuth service health
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to introspect a dummy token to check service availability
	_, err := gw.authClient.IntrospectToken(ctx, "health-check")

	status := "healthy"
	statusCode := http.StatusOK

	if err != nil {
		// Service is reachable but token is invalid (expected)
		if !contains(err.Error(), "invalid") && !contains(err.Error(), "unauthorized") {
			status = "unhealthy"
			statusCode = http.StatusServiceUnavailable
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"status":    status,
		"service":   "api-gateway",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && contains(s[1:], substr)
}

// SetupRoutes sets up all API Gateway routes
func (gw *APIGateway) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", gw.HandleHealthCheck).Methods("GET")

	// Authentication routes
	r.HandleFunc("/api/auth/login", gw.HandleLogin).Methods("POST")
	r.HandleFunc("/api/auth/register", gw.HandleRegister).Methods("POST")
	r.HandleFunc("/api/auth/change-password", gw.HandleChangePassword).Methods("POST")
	r.HandleFunc("/api/auth/introspect", gw.HandleTokenIntrospect).Methods("POST")

	// User routes
	r.HandleFunc("/api/users/me", gw.HandleGetCurrentUser).Methods("GET")

	// Admin routes
	r.HandleFunc("/api/admin/roles", gw.HandleCreateRole).Methods("POST")
	r.HandleFunc("/api/admin/roles", gw.HandleListRoles).Methods("GET")
	r.HandleFunc("/api/admin/users/{userId}/roles", gw.HandleAssignRoleToUser).Methods("POST")

	// Authorization routes
	r.HandleFunc("/api/authz/check", gw.HandleCheckPermission).Methods("POST")

	return r
}

func main() {
	// Get server configuration from environment
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	// Initialize API Gateway
	gateway := NewAPIGateway()
	router := gateway.SetupRoutes()

	// Start server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("API Gateway starting on port %s", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)
	log.Fatal(server.ListenAndServe())
}
