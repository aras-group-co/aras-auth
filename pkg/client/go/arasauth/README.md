# ArasAuth Go SDK

A comprehensive Go client library for the ArasAuth authentication and authorization service. This SDK provides complete access to all authentication, user management, group management, role-based access control (RBAC), and authorization features.

## Features

- **Complete API Coverage**: 42/42 endpoints (100% coverage)
- **Authentication**: Login, register, password management, token introspection
- **User Management**: Full CRUD operations for users
- **Group Management**: Complete group lifecycle management
- **Role-Based Access Control**: Role and permission management
- **Authorization**: Permission checking and role assignments
- **Type-Safe**: Strongly typed Go structs for all requests and responses
- **Context Support**: Full context.Context support for timeouts and cancellation
- **Error Handling**: Comprehensive error handling with detailed error messages

## Installation

```bash
go get github.com/aras-services/aras-auth/pkg/client/go/arasauth
```

## Quick Start

### Basic Authentication

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/aras-services/aras-auth/pkg/client/go/arasauth"
)

func main() {
    // Initialize client
    client := arasauth.NewClient("http://localhost:7600")
    
    // Login
    authResp, err := client.Login(context.Background(), "user@example.com", "password")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Access token: %s\n", authResp.AccessToken)
    
    // Get current user
    user, err := client.GetCurrentUser(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User: %s %s\n", user.FirstName, user.LastName)
}
```

## API Gateway Usage

This SDK is specifically designed for use in API Gateway services where you need to forward requests from frontend applications to the Auth service.

### Example: API Gateway Handler

```go
package main

import (
    "context"
    "encoding/json"
    "net/http"
    
    "github.com/aras-services/aras-auth/pkg/client/go/arasauth"
    "github.com/gorilla/mux"
)

type APIGateway struct {
    authClient *arasauth.Client
}

func NewAPIGateway(authServiceURL string) *APIGateway {
    return &APIGateway{
        authClient: arasauth.NewClient(authServiceURL),
    }
}

// Handle user login from frontend
func (gw *APIGateway) HandleLogin(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    // Forward to auth service
    authResp, err := gw.authClient.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        http.Error(w, "Authentication failed", http.StatusUnauthorized)
        return
    }
    
    // Return response to frontend
    json.NewEncoder(w).Encode(authResp)
}

// Handle password change from frontend
func (gw *APIGateway) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
    // Extract token from Authorization header
    token := extractTokenFromHeader(r)
    gw.authClient.SetToken(token)
    
    var req struct {
        CurrentPassword string `json:"current_password"`
        NewPassword     string `json:"new_password"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    // Forward to auth service
    err := gw.authClient.ChangePassword(r.Context(), req.CurrentPassword, req.NewPassword)
    if err != nil {
        http.Error(w, "Password change failed", http.StatusBadRequest)
        return
    }
    
    w.WriteHeader(http.StatusOK)
}

// Handle admin role management from frontend
func (gw *APIGateway) HandleCreateRole(w http.ResponseWriter, r *http.Request) {
    // Extract admin token
    token := extractTokenFromHeader(r)
    gw.authClient.SetToken(token)
    
    var req struct {
        Name        string `json:"name"`
        Description string `json:"description"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    // Forward to auth service
    role, err := gw.authClient.CreateRole(r.Context(), req.Name, req.Description)
    if err != nil {
        http.Error(w, "Role creation failed", http.StatusBadRequest)
        return
    }
    
    json.NewEncoder(w).Encode(role)
}

func extractTokenFromHeader(r *http.Request) string {
    auth := r.Header.Get("Authorization")
    if len(auth) > 7 && auth[:7] == "Bearer " {
        return auth[7:]
    }
    return ""
}
```

## Complete API Reference

### Authentication APIs

```go
// Basic authentication
authResp, err := client.Login(ctx, "user@example.com", "password")
user, err := client.Register(ctx, "user@example.com", "password", "John", "Doe")
err = client.Logout(ctx, refreshToken)
authResp, err = client.RefreshToken(ctx, refreshToken)

// Advanced authentication
err = client.ChangePassword(ctx, currentPassword, newPassword)
err = client.ForgotPassword(ctx, "user@example.com")
err = client.ResetPassword(ctx, resetToken, newPassword)
err = client.VerifyEmail(ctx, userID)

// Token introspection
introspection, err := client.IntrospectToken(ctx, token)
```

### User Management APIs

```go
// User operations
user, err := client.GetCurrentUser(ctx)
users, err := client.ListUsers(ctx, page, limit)
user, err := client.GetUser(ctx, userID)
user, err := client.UpdateUser(ctx, userID, &UpdateUserRequest{...})
err = client.DeleteUser(ctx, userID)
```

### Group Management APIs

```go
// Group operations
group, err := client.CreateGroup(ctx, "Developers", "Development team")
groups, err := client.ListGroups(ctx, page, limit)
group, err := client.GetGroup(ctx, groupID)
group, err := client.UpdateGroup(ctx, groupID, &UpdateGroupRequest{...})
err = client.DeleteGroup(ctx, groupID)

// Group membership
err = client.AddMember(ctx, groupID, userID)
err = client.RemoveMember(ctx, groupID, userID)
members, err := client.GetMembers(ctx, groupID)
```

### Role Management APIs

```go
// Role operations
role, err := client.CreateRole(ctx, "admin", "Administrator role")
roles, err := client.ListRoles(ctx, page, limit)
role, err := client.GetRole(ctx, roleID)
role, err := client.UpdateRole(ctx, roleID, &UpdateRoleRequest{...})
err = client.DeleteRole(ctx, roleID)

// Role permissions
err = client.AssignPermissionToRole(ctx, roleID, permissionID)
err = client.RemovePermissionFromRole(ctx, roleID, permissionID)
permissions, err := client.GetRolePermissions(ctx, roleID)
```

### Permission Management APIs

```go
// Permission operations
permission, err := client.CreatePermission(ctx, "users", "read", "Read user information")
permissions, err := client.ListPermissions(ctx, page, limit)
permission, err := client.GetPermission(ctx, permissionID)
permission, err := client.UpdatePermission(ctx, permissionID, &UpdatePermissionRequest{...})
err = client.DeletePermission(ctx, permissionID)
```

### Authorization APIs

```go
// Role assignments
err = client.AssignRoleToUser(ctx, userID, roleID)
err = client.RemoveRoleFromUser(ctx, userID, roleID)
roles, err := client.GetUserRoles(ctx, userID)

err = client.AssignRoleToGroup(ctx, groupID, roleID)
err = client.RemoveRoleFromGroup(ctx, groupID, roleID)
roles, err = client.GetGroupRoles(ctx, groupID)

// Permission checking
hasPermission, err := client.CheckPermission(ctx, userID, "users", "read")
```

## Data Models

### Core Models

```go
type User struct {
    ID            string `json:"id"`
    Email         string `json:"email"`
    FirstName     string `json:"first_name"`
    LastName      string `json:"last_name"`
    Status        string `json:"status"`
    EmailVerified bool   `json:"email_verified"`
    CreatedAt     string `json:"created_at"`
    UpdatedAt     string `json:"updated_at"`
}

type Group struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type Role struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type Permission struct {
    ID          string `json:"id"`
    Resource    string `json:"resource"`
    Action      string `json:"action"`
    Description string `json:"description"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}
```

### Request/Response Models

```go
// Authentication
type AuthResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
    TokenType    string `json:"token_type"`
    User         *User  `json:"user"`
}

type TokenIntrospection struct {
    Active    bool   `json:"active"`
    UserID    string `json:"user_id,omitempty"`
    Email     string `json:"email,omitempty"`
    ExpiresAt int64  `json:"exp,omitempty"`
}

// Pagination
type ListUsersResponse struct {
    Users []*User `json:"users"`
    Total int     `json:"total"`
    Page  int     `json:"page"`
    Limit int     `json:"limit"`
}

type ListGroupsResponse struct {
    Groups []*Group `json:"groups"`
    Total  int      `json:"total"`
    Page   int      `json:"page"`
    Limit  int      `json:"limit"`
}

type ListRolesResponse struct {
    Roles []*Role `json:"roles"`
    Total int     `json:"total"`
    Page  int     `json:"page"`
    Limit int     `json:"limit"`
}

type ListPermissionsResponse struct {
    Permissions []*Permission `json:"permissions"`
    Total       int           `json:"total"`
    Page        int           `json:"page"`
    Limit       int           `json:"limit"`
}
```

## Error Handling

The SDK provides comprehensive error handling:

```go
user, err := client.GetUser(ctx, userID)
if err != nil {
    // Handle different types of errors
    switch {
    case strings.Contains(err.Error(), "not found"):
        // User not found
    case strings.Contains(err.Error(), "unauthorized"):
        // Authentication required
    case strings.Contains(err.Error(), "forbidden"):
        // Insufficient permissions
    default:
        // Other errors
    }
}
```

## Context Support

All methods support context for timeouts and cancellation:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := client.GetUser(ctx, userID)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(2 * time.Second)
    cancel() // Cancel the request
}()
user, err := client.GetUser(ctx, userID)
```

## Configuration

### Client Configuration

```go
client := arasauth.NewClient("http://localhost:7600")

// Set authentication token for protected endpoints
client.SetToken("your-access-token")

// Custom HTTP client (optional)
client := arasauth.NewClient("http://localhost:7600")
// The client uses a default 30-second timeout
```

## Best Practices

### 1. Token Management

```go
// Always set token before making protected requests
client.SetToken(accessToken)

// Check token validity before making requests
introspection, err := client.IntrospectToken(ctx, token)
if err != nil || !introspection.Active {
    // Token is invalid, redirect to login
}
```

### 2. Error Handling

```go
// Always handle errors appropriately
user, err := client.GetUser(ctx, userID)
if err != nil {
    log.Printf("Failed to get user: %v", err)
    return
}
```

### 3. Pagination

```go
// Use pagination for list operations
users, err := client.ListUsers(ctx, 1, 20) // page 1, limit 20
if err != nil {
    return err
}

fmt.Printf("Total users: %d\n", users.Total)
for _, user := range users.Users {
    fmt.Printf("User: %s\n", user.Email)
}
```

### 4. Context Usage

```go
// Always use context for request cancellation and timeouts
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Pass context to all SDK calls
user, err := client.GetUser(ctx, userID)
```

## License

This SDK is part of the ArasAuth project and is licensed under the MIT License.

## Integration in Production

### Service Configuration

When integrating ArasAuth in your service, use environment variables for the Auth service URL:

```go
package main

import (
    "os"
    "github.com/aras-services/aras-auth/pkg/client/go/arasauth"
)

func main() {
    // Get Auth service URL from environment
    authURL := os.Getenv("ARAS_AUTH_URL")
    if authURL == "" {
        authURL = "http://localhost:7600" // fallback for local dev
    }
    
    client := arasauth.NewClient(authURL)
    // ... your service logic
}
```

### Docker Compose Integration

Example `docker-compose.yml` for a service using ArasAuth:

```yaml
version: '3.8'

services:
  # ArasAuth service
  aras_auth:
    image: ghcr.io/aras-group-co/aras-auth:v1.3.0
    environment:
      SERVER_HOST: 0.0.0.0
      SERVER_PORT: 7600
      DB_HOST: postgres
      DB_USER: ${POSTGRES_USER:-postgres}
      DB_PASSWORD: ${POSTGRES_PASSWORD:-postgres}
      DB_NAME: ${POSTGRES_DB:-myproduct}
      JWT_SECRET_KEY: ${JWT_SECRET:-change-me-32-chars-min}
      ADMIN_EMAIL: admin@myproduct.com
    networks:
      - app_network
  
  # Your service
  my_service:
    image: mycompany/my-service:latest
    environment:
      # Point to ArasAuth service
      ARAS_AUTH_URL: http://aras_auth:7600
    depends_on:
      - aras_auth
    networks:
      - app_network

networks:
  app_network:
    driver: bridge
```

### Environment Variables

Required environment variables for services using ArasAuth SDK:

| Variable | Description | Example |
|----------|-------------|---------|
| `ARAS_AUTH_URL` | ArasAuth service URL | `http://aras_auth:7600` |

### Best Practices

1. **Always use environment variables** for service URLs
2. **Handle errors gracefully** with proper retry logic
3. **Use context with timeouts** for all API calls
4. **Cache user permissions** to reduce API calls
5. **Implement circuit breaker** for resilience

### Example: Production-Ready Service

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"
    
    "github.com/aras-services/aras-auth/pkg/client/go/arasauth"
)

type MyService struct {
    authClient *arasauth.Client
}

func NewMyService() *MyService {
    authURL := os.Getenv("ARAS_AUTH_URL")
    if authURL == "" {
        log.Fatal("ARAS_AUTH_URL environment variable is required")
    }
    
    return &MyService{
        authClient: arasauth.NewClient(authURL),
    }
}

func (s *MyService) AuthenticateRequest(token string) (*arasauth.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    s.authClient.SetToken(token)
    
    user, err := s.authClient.GetCurrentUser(ctx)
    if err != nil {
        return nil, fmt.Errorf("authentication failed: %w", err)
    }
    
    return user, nil
}

func (s *MyService) CheckPermission(userID, resource, action string) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    hasPermission, err := s.authClient.CheckPermission(ctx, userID, resource, action)
    if err != nil {
        return false, fmt.Errorf("permission check failed: %w", err)
    }
    
    return hasPermission, nil
}

func main() {
    service := NewMyService()
    
    // Your service logic here
    log.Println("Service started with ArasAuth integration")
}
```

## Support

For issues and questions:
- GitHub Issues: [aras-auth/issues](https://github.com/aras-services/aras-auth/issues)
- Documentation: [aras-auth/docs](https://github.com/aras-services/aras-auth/wiki)

