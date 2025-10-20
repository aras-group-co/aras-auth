# ArasAuth Integration Guide

A comprehensive guide for integrating ArasAuth authentication and authorization service into your applications and microservices.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Environment Configuration](#environment-configuration)
3. [Docker Compose Templates](#docker-compose-templates)
4. [Complete Examples](#complete-examples)
5. [Testing](#testing)
6. [Troubleshooting](#troubleshooting)
7. [Performance Tips](#performance-tips)
8. [Security Best Practices](#security-best-practices)

## Quick Start

### 1. Install the SDK

```bash
go get github.com/aras-services/aras-auth/pkg/client/go/arasauth
```

### 2. Basic Integration

```go
package main

import (
    "context"
    "log"
    "os"
    
    "github.com/aras-services/aras-auth/pkg/client/go/arasauth"
)

func main() {
    // Get Auth service URL from environment
    authURL := os.Getenv("ARAS_AUTH_URL")
    if authURL == "" {
        authURL = "http://localhost:7600" // fallback for local dev
    }
    
    // Initialize client
    client := arasauth.NewClient(authURL)
    
    // Test connection
    ctx := context.Background()
    _, err := client.IntrospectToken(ctx, "test-token")
    if err != nil {
        log.Printf("Auth service connection test: %v", err)
    } else {
        log.Println("Auth service is reachable")
    }
}
```

### 3. Deploy ArasAuth Service

```bash
# Set version
export ARAS_AUTH_VERSION=v1.3.0

# Deploy with docker-compose
docker compose up -d
```

## Environment Configuration

### Required Environment Variables

| Variable | Description | Example | Required |
|----------|-------------|---------|----------|
| `ARAS_AUTH_URL` | ArasAuth service URL | `http://aras_auth:7600` | ✅ Yes |

### Optional Environment Variables

| Variable | Description | Default | Use Case |
|----------|-------------|---------|----------|
| `ARAS_AUTH_TIMEOUT` | Request timeout in seconds | `30` | Custom timeout |
| `ARAS_AUTH_RETRY_COUNT` | Number of retries | `3` | Resilience |
| `ARAS_AUTH_CACHE_TTL` | Cache TTL in seconds | `300` | Performance |

### Example .env File

```bash
# ArasAuth Configuration
ARAS_AUTH_URL=http://aras_auth:7600
ARAS_AUTH_TIMEOUT=10
ARAS_AUTH_RETRY_COUNT=3
ARAS_AUTH_CACHE_TTL=300

# Your Service Configuration
SERVICE_PORT=8080
SERVICE_NAME=my-service
```

## Docker Compose Templates

### Template 1: Basic Integration

```yaml
version: '3.8'

services:
  # ArasAuth Service
  aras_auth:
    image: ghcr.io/aras-group-co/aras-auth:${ARAS_AUTH_VERSION}
    container_name: aras_auth_service
    environment:
      SERVER_HOST: 0.0.0.0
      SERVER_PORT: 7600
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${POSTGRES_USER:-postgres}
      DB_PASSWORD: ${POSTGRES_PASSWORD:-postgres}
      DB_NAME: ${POSTGRES_DB:-myproduct}
      DB_SSL_MODE: disable
      JWT_SECRET_KEY: ${JWT_SECRET:-change-me-32-chars-min}
      JWT_ACCESS_EXPIRY: 15m
      JWT_REFRESH_EXPIRY: 168h
      ADMIN_EMAIL: ${ADMIN_EMAIL:-admin@myproduct.com}
      ADMIN_PASSWORD: ${ADMIN_PASSWORD:-admin123}
    ports:
      - "7600:7600"
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - app_network

  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: postgres_db
    environment:
      POSTGRES_USER: ${POSTGRES_USER:-postgres}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-postgres}
      POSTGRES_DB: ${POSTGRES_DB:-myproduct}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-postgres}"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app_network

  # Your Service
  my_service:
    image: mycompany/my-service:latest
    container_name: my_service
    environment:
      ARAS_AUTH_URL: http://aras_auth:7600
      SERVICE_PORT: 8080
    ports:
      - "8080:8080"
    depends_on:
      - aras_auth
    networks:
      - app_network

volumes:
  postgres_data:

networks:
  app_network:
    driver: bridge
```

### Template 2: Microservices Architecture

```yaml
version: '3.8'

services:
  # ArasAuth Service
  aras_auth:
    image: ghcr.io/aras-group-co/aras-auth:${ARAS_AUTH_VERSION}
    environment:
      SERVER_HOST: 0.0.0.0
      SERVER_PORT: 7600
      DB_HOST: postgres
      DB_USER: ${POSTGRES_USER:-postgres}
      DB_PASSWORD: ${POSTGRES_PASSWORD:-postgres}
      DB_NAME: ${POSTGRES_DB:-myproduct}
      JWT_SECRET_KEY: ${JWT_SECRET:-change-me-32-chars-min}
    networks:
      - backend_network

  # API Gateway
  api_gateway:
    image: mycompany/api-gateway:latest
    environment:
      ARAS_AUTH_URL: http://aras_auth:7600
    ports:
      - "80:8080"
    depends_on:
      - aras_auth
    networks:
      - backend_network

  # User Service
  user_service:
    image: mycompany/user-service:latest
    environment:
      ARAS_AUTH_URL: http://aras_auth:7600
    depends_on:
      - aras_auth
    networks:
      - backend_network

  # Order Service
  order_service:
    image: mycompany/order-service:latest
    environment:
      ARAS_AUTH_URL: http://aras_auth:7600
    depends_on:
      - aras_auth
    networks:
      - backend_network

  # Database
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: ${POSTGRES_USER:-postgres}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-postgres}
      POSTGRES_DB: ${POSTGRES_DB:-myproduct}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - backend_network

volumes:
  postgres_data:

networks:
  backend_network:
    driver: bridge
```

## Complete Examples

### Example 1: Authentication Middleware

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
    
    "github.com/aras-services/aras-auth/pkg/client/go/arasauth"
    "github.com/gorilla/mux"
)

type AuthMiddleware struct {
    authClient *arasauth.Client
}

func NewAuthMiddleware() *AuthMiddleware {
    authURL := os.Getenv("ARAS_AUTH_URL")
    if authURL == "" {
        log.Fatal("ARAS_AUTH_URL environment variable is required")
    }
    
    return &AuthMiddleware{
        authClient: arasauth.NewClient(authURL),
    }
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract token from Authorization header
        auth := r.Header.Get("Authorization")
        if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        
        token := strings.TrimPrefix(auth, "Bearer ")
        
        // Validate token with ArasAuth
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        introspection, err := m.authClient.IntrospectToken(ctx, token)
        if err != nil || !introspection.Active {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // Add user info to request context
        ctx = context.WithValue(ctx, "user_id", introspection.UserID)
        ctx = context.WithValue(ctx, "user_email", introspection.Email)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

func (m *AuthMiddleware) RequirePermission(resource, action string) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            userID := r.Context().Value("user_id").(string)
            
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()
            
            m.authClient.SetToken(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
            
            hasPermission, err := m.authClient.CheckPermission(ctx, userID, resource, action)
            if err != nil || !hasPermission {
                http.Error(w, "Insufficient permissions", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        }
    }
}

func main() {
    authMiddleware := NewAuthMiddleware()
    
    r := mux.NewRouter()
    
    // Public routes
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
    }).Methods("GET")
    
    // Protected routes
    r.HandleFunc("/api/users", authMiddleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("user_id").(string)
        json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
    })).Methods("GET")
    
    r.HandleFunc("/api/admin/users", authMiddleware.RequireAuth(
        authMiddleware.RequirePermission("users", "admin")(
            func(w http.ResponseWriter, r *http.Request) {
                json.NewEncoder(w).Encode(map[string]string{"message": "Admin access granted"})
            },
        ),
    )).Methods("GET")
    
    server := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }
    
    log.Println("Server starting on :8080")
    log.Fatal(server.ListenAndServe())
}
```

### Example 2: Service with Caching

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "sync"
    "time"
    
    "github.com/aras-services/aras-auth/pkg/client/go/arasauth"
)

type CachedAuthService struct {
    authClient *arasauth.Client
    cache      map[string]*CacheEntry
    mutex      sync.RWMutex
    ttl        time.Duration
}

type CacheEntry struct {
    User        *arasauth.User
    Permissions map[string]bool
    ExpiresAt   time.Time
}

func NewCachedAuthService() *CachedAuthService {
    authURL := os.Getenv("ARAS_AUTH_URL")
    if authURL == "" {
        log.Fatal("ARAS_AUTH_URL environment variable is required")
    }
    
    ttl := 5 * time.Minute // Default cache TTL
    if ttlStr := os.Getenv("ARAS_AUTH_CACHE_TTL"); ttlStr != "" {
        if parsedTTL, err := time.ParseDuration(ttlStr + "s"); err == nil {
            ttl = parsedTTL
        }
    }
    
    return &CachedAuthService{
        authClient: arasauth.NewClient(authURL),
        cache:      make(map[string]*CacheEntry),
        ttl:        ttl,
    }
}

func (s *CachedAuthService) GetUser(token string) (*arasauth.User, error) {
    s.mutex.RLock()
    entry, exists := s.cache[token]
    s.mutex.RUnlock()
    
    if exists && time.Now().Before(entry.ExpiresAt) {
        return entry.User, nil
    }
    
    // Cache miss or expired, fetch from auth service
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    s.authClient.SetToken(token)
    user, err := s.authClient.GetCurrentUser(ctx)
    if err != nil {
        return nil, err
    }
    
    // Update cache
    s.mutex.Lock()
    s.cache[token] = &CacheEntry{
        User:      user,
        ExpiresAt: time.Now().Add(s.ttl),
    }
    s.mutex.Unlock()
    
    return user, nil
}

func (s *CachedAuthService) CheckPermission(token, userID, resource, action string) (bool, error) {
    s.mutex.RLock()
    entry, exists := s.cache[token]
    s.mutex.RUnlock()
    
    if exists && time.Now().Before(entry.ExpiresAt) {
        if permissions, ok := entry.Permissions[resource+":"+action]; ok {
            return permissions, nil
        }
    }
    
    // Cache miss or expired, fetch from auth service
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    s.authClient.SetToken(token)
    hasPermission, err := s.authClient.CheckPermission(ctx, userID, resource, action)
    if err != nil {
        return false, err
    }
    
    // Update cache
    s.mutex.Lock()
    if entry == nil {
        entry = &CacheEntry{
            Permissions: make(map[string]bool),
            ExpiresAt:   time.Now().Add(s.ttl),
        }
    }
    entry.Permissions[resource+":"+action] = hasPermission
    s.cache[token] = entry
    s.mutex.Unlock()
    
    return hasPermission, nil
}

func (s *CachedAuthService) ClearCache() {
    s.mutex.Lock()
    s.cache = make(map[string]*CacheEntry)
    s.mutex.Unlock()
}

func main() {
    authService := NewCachedAuthService()
    
    // Example usage
    token := "your-jwt-token"
    
    user, err := authService.GetUser(token)
    if err != nil {
        log.Printf("Failed to get user: %v", err)
        return
    }
    
    fmt.Printf("User: %s %s\n", user.FirstName, user.LastName)
    
    hasPermission, err := authService.CheckPermission(token, user.ID, "users", "read")
    if err != nil {
        log.Printf("Failed to check permission: %v", err)
        return
    }
    
    fmt.Printf("Has permission: %v\n", hasPermission)
}
```

### Example 3: Circuit Breaker Pattern

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "sync"
    "time"
    
    "github.com/aras-services/aras-auth/pkg/client/go/arasauth"
)

type CircuitBreakerState int

const (
    StateClosed CircuitBreakerState = iota
    StateOpen
    StateHalfOpen
)

type CircuitBreaker struct {
    state         CircuitBreakerState
    failureCount  int
    lastFailTime  time.Time
    mutex         sync.RWMutex
    maxFailures   int
    timeout       time.Duration
}

type ResilientAuthService struct {
    authClient     *arasauth.Client
    circuitBreaker *CircuitBreaker
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        state:       StateClosed,
        maxFailures: maxFailures,
        timeout:     timeout,
    }
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    if cb.state == StateOpen {
        if time.Since(cb.lastFailTime) > cb.timeout {
            cb.state = StateHalfOpen
        } else {
            return fmt.Errorf("circuit breaker is open")
        }
    }
    
    err := fn()
    
    if err != nil {
        cb.failureCount++
        cb.lastFailTime = time.Now()
        
        if cb.failureCount >= cb.maxFailures {
            cb.state = StateOpen
        }
        
        return err
    }
    
    // Success
    cb.failureCount = 0
    cb.state = StateClosed
    return nil
}

func NewResilientAuthService() *ResilientAuthService {
    authURL := os.Getenv("ARAS_AUTH_URL")
    if authURL == "" {
        log.Fatal("ARAS_AUTH_URL environment variable is required")
    }
    
    return &ResilientAuthService{
        authClient:     arasauth.NewClient(authURL),
        circuitBreaker: NewCircuitBreaker(5, 30*time.Second),
    }
}

func (s *ResilientAuthService) GetUserWithRetry(token string) (*arasauth.User, error) {
    var user *arasauth.User
    var err error
    
    circuitErr := s.circuitBreaker.Execute(func() error {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        s.authClient.SetToken(token)
        user, err = s.authClient.GetCurrentUser(ctx)
        return err
    })
    
    if circuitErr != nil {
        return nil, circuitErr
    }
    
    return user, nil
}

func main() {
    authService := NewResilientAuthService()
    
    // Example usage with retry
    token := "your-jwt-token"
    
    user, err := authService.GetUserWithRetry(token)
    if err != nil {
        log.Printf("Failed to get user: %v", err)
        return
    }
    
    fmt.Printf("User: %s %s\n", user.FirstName, user.LastName)
}
```

## Testing

### Unit Tests

```go
package main

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/aras-services/aras-auth/pkg/client/go/arasauth"
)

func TestAuthMiddleware(t *testing.T) {
    // Mock ArasAuth service
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/introspect" {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{"active": true, "user_id": "123", "email": "test@example.com"}`))
        }
    }))
    defer server.Close()
    
    // Create client with mock server
    client := arasauth.NewClient(server.URL)
    middleware := &AuthMiddleware{authClient: client}
    
    // Test protected endpoint
    req := httptest.NewRequest("GET", "/api/users", nil)
    req.Header.Set("Authorization", "Bearer valid-token")
    
    w := httptest.NewRecorder()
    
    middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("user_id").(string)
        if userID != "123" {
            t.Errorf("Expected user_id 123, got %s", userID)
        }
    })(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

### Integration Tests

```go
package main

import (
    "context"
    "os"
    "testing"
    "time"
    
    "github.com/aras-services/aras-auth/pkg/client/go/arasauth"
)

func TestIntegration(t *testing.T) {
    // Skip if not in integration test environment
    if os.Getenv("INTEGRATION_TESTS") != "true" {
        t.Skip("Skipping integration test")
    }
    
    authURL := os.Getenv("ARAS_AUTH_URL")
    if authURL == "" {
        authURL = "http://localhost:7600"
    }
    
    client := arasauth.NewClient(authURL)
    
    // Test service health
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    _, err := client.IntrospectToken(ctx, "test-token")
    if err == nil {
        t.Error("Expected error for invalid token")
    }
    
    // Test login (requires valid credentials)
    authResp, err := client.Login(ctx, "admin@aras-services.com", "admin123")
    if err != nil {
        t.Logf("Login test skipped: %v", err)
        return
    }
    
    if authResp.AccessToken == "" {
        t.Error("Expected access token")
    }
    
    // Test authenticated request
    client.SetToken(authResp.AccessToken)
    user, err := client.GetCurrentUser(ctx)
    if err != nil {
        t.Errorf("Failed to get current user: %v", err)
    }
    
    if user.Email != "admin@aras-services.com" {
        t.Errorf("Expected admin email, got %s", user.Email)
    }
}
```

## Troubleshooting

### Common Issues

#### 1. Connection Refused

**Error**: `dial tcp: connection refused`

**Solutions**:
- Check if ArasAuth service is running: `docker ps | grep aras_auth`
- Verify service URL: `echo $ARAS_AUTH_URL`
- Check network connectivity: `curl http://aras_auth:7600/health`

#### 2. Authentication Failed

**Error**: `authentication failed` or `invalid token`

**Solutions**:
- Verify JWT secret is consistent across services
- Check token expiration
- Ensure token format: `Bearer <token>`

#### 3. Permission Denied

**Error**: `insufficient permissions`

**Solutions**:
- Check user roles and permissions
- Verify resource and action names
- Ensure user has required permissions

#### 4. Database Connection Issues

**Error**: `database connection failed`

**Solutions**:
- Check PostgreSQL service status
- Verify database credentials
- Check network connectivity between services

### Debug Mode

Enable debug logging:

```go
import "log"

func main() {
    // Enable debug logging
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    // Your service code
}
```

### Health Checks

```go
func healthCheck(authURL string) error {
    client := arasauth.NewClient(authURL)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    _, err := client.IntrospectToken(ctx, "health-check")
    if err != nil {
        return fmt.Errorf("auth service health check failed: %w", err)
    }
    
    return nil
}
```

## Performance Tips

### 1. Connection Pooling

```go
import (
    "net/http"
    "time"
)

func createHTTPClient() *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
        },
        Timeout: 30 * time.Second,
    }
}
```

### 2. Caching Strategy

```go
type CacheConfig struct {
    TTL           time.Duration
    MaxSize       int
    CleanupInterval time.Duration
}

func NewCacheConfig() *CacheConfig {
    return &CacheConfig{
        TTL:            5 * time.Minute,
        MaxSize:        1000,
        CleanupInterval: 1 * time.Minute,
    }
}
```

### 3. Batch Operations

```go
func (s *AuthService) BatchCheckPermissions(userID string, checks []PermissionCheck) ([]bool, error) {
    results := make([]bool, len(checks))
    
    for i, check := range checks {
        hasPermission, err := s.CheckPermission(userID, check.Resource, check.Action)
        if err != nil {
            return nil, err
        }
        results[i] = hasPermission
    }
    
    return results, nil
}
```

### 4. Async Operations

```go
func (s *AuthService) AsyncGetUser(token string) <-chan *arasauth.User {
    result := make(chan *arasauth.User, 1)
    
    go func() {
        defer close(result)
        
        user, err := s.GetUser(token)
        if err != nil {
            log.Printf("Failed to get user: %v", err)
            return
        }
        
        result <- user
    }()
    
    return result
}
```

## Security Best Practices

### 1. Token Validation

```go
func validateToken(token string) error {
    if len(token) < 10 {
        return fmt.Errorf("token too short")
    }
    
    // Additional validation logic
    return nil
}
```

### 2. Secure Headers

```go
func addSecurityHeaders(w http.ResponseWriter) {
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "DENY")
    w.Header().Set("X-XSS-Protection", "1; mode=block")
    w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
}
```

### 3. Rate Limiting

```go
import "golang.org/x/time/rate"

type RateLimiter struct {
    limiter *rate.Limiter
}

func NewRateLimiter(rps int) *RateLimiter {
    return &RateLimiter{
        limiter: rate.NewLimiter(rate.Limit(rps), 1),
    }
}

func (rl *RateLimiter) Allow() bool {
    return rl.limiter.Allow()
}
```

### 4. Input Validation

```go
import "github.com/go-playground/validator/v10"

type UserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func validateRequest(req *UserRequest) error {
    validate := validator.New()
    return validate.Struct(req)
}
```

### 5. Logging Security

```go
func sanitizeLogData(data interface{}) interface{} {
    // Remove sensitive information from logs
    switch v := data.(type) {
    case string:
        if strings.Contains(v, "password") || strings.Contains(v, "token") {
            return "[REDACTED]"
        }
    }
    return data
}
```

## Support

For additional help:

- **GitHub Issues**: [aras-auth/issues](https://github.com/aras-services/aras-auth/issues)
- **Documentation**: [aras-auth/docs](https://github.com/aras-services/aras-auth/wiki)
- **Examples**: [pkg/client/go/examples](https://github.com/aras-services/aras-auth/tree/main/pkg/client/go/examples)

## License

This integration guide is part of the ArasAuth project and is licensed under the MIT License.
