# ArasAuth - Authentication & Authorization Service

A production-ready, reusable authentication and authorization microservice built with Go, following Clean Architecture principles and inspired by the ORY stack. This service provides comprehensive user management, group management, role-based access control (RBAC), and JWT-based authentication.

## 🚀 Features

### Core Features
- **User Management**: Registration, login, profile management
- **Group Management**: Create and manage user groups
- **Role-Based Access Control (RBAC)**: Fine-grained permissions system
- **JWT Authentication**: Access tokens with refresh token rotation
- **Provider Architecture**: Pluggable identity providers (currently local, extensible for LDAP, OAuth, etc.)
- **RESTful API**: Comprehensive REST API for all operations
- **Multi-language SDKs**: Go and Python client libraries
- **Docker Support**: Production-ready containerization

### Security Features
- Password hashing with bcrypt (cost 12)
- JWT with HS256 signing
- Refresh token rotation
- Rate limiting on auth endpoints
- CORS configuration
- Secure headers middleware
- Input validation and sanitization

## 🏗️ Architecture

### Clean Architecture Layers
```
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── domain/            # Business entities & interfaces (core)
│   ├── usecase/           # Business logic
│   ├── delivery/          # HTTP handlers (REST API)
│   ├── repository/        # Data access implementations
│   ├── provider/          # Identity providers
│   └── middleware/        # Auth middleware
├── pkg/                   # Public libraries (SDK)
├── config/                # Configuration management
├── migrations/            # PostgreSQL migrations
└── scripts/               # Utility scripts
```

### Core Components (ORY-inspired)

1. **Identity Service** (Kratos-like)
   - User registration, login, password management
   - Email verification, password reset flows
   - Local provider implementation with PostgreSQL

2. **Token Service** (Hydra-like - simplified JWT version)
   - JWT access token generation/validation
   - Refresh token management with session storage
   - Token introspection endpoint

3. **Authorization Service** (Keto-like)
   - RBAC implementation (roles, permissions)
   - Permission checking endpoints
   - Group-based access control

4. **Gateway Integration** (Oathkeeper-like)
   - Middleware for other services
   - Token validation
   - Permission enforcement

## 📋 Prerequisites

- Go 1.22+
- PostgreSQL 15+
- Docker & Docker Compose (for containerized deployment)

## 🚀 Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/aras-services/aras-auth.git
   cd aras-auth
   ```

2. **Start the services**
   ```bash
   docker-compose up -d
   ```

3. **Run migrations**
   ```bash
   docker-compose exec aras_auth ./main migrate up
   ```

4. **Access the service**
   - API: http://localhost:8080
   - Health check: http://localhost:8080/health

### Manual Setup

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL**
   ```bash
   createdb aras_auth
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run migrations**
   ```bash
   go run ./cmd/migrate up
   ```

5. **Start the service**
   ```bash
   go run ./cmd/server
   ```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | Server host | `0.0.0.0` |
| `SERVER_PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `aras_auth` |
| `DB_SSL_MODE` | SSL mode | `disable` |
| `JWT_SECRET_KEY` | JWT secret key | `your-secret-key` |
| `JWT_ACCESS_EXPIRY` | Access token expiry | `15m` |
| `JWT_REFRESH_EXPIRY` | Refresh token expiry | `7d` |
| `ADMIN_EMAIL` | Admin email | `admin@aras-services.com` |
| `ADMIN_PASSWORD` | Admin password | `admin123` |

### Configuration File

You can also use a YAML configuration file (`config/config.yaml`):

```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  name: "aras_auth"
  ssl_mode: "disable"

jwt:
  secret_key: "your-secret-key-change-in-production"
  access_expiry: "15m"
  refresh_expiry: "7d"

admin:
  email: "admin@aras-services.com"
  password: "admin123"
```

## 📚 API Documentation

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900,
    "token_type": "Bearer",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "status": "active",
      "email_verified": false,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Logout
```http
POST /api/v1/auth/logout
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### User Management Endpoints

#### Get Current User
```http
GET /api/v1/users/me
Authorization: Bearer <access_token>
```

#### List Users
```http
GET /api/v1/users?page=1&limit=20
Authorization: Bearer <access_token>
```

#### Get User by ID
```http
GET /api/v1/users/{id}
Authorization: Bearer <access_token>
```

#### Update User
```http
PUT /api/v1/users/{id}
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith"
}
```

#### Delete User
```http
DELETE /api/v1/users/{id}
Authorization: Bearer <access_token>
```

### Group Management Endpoints

#### Create Group
```http
POST /api/v1/groups
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Developers",
  "description": "Development team"
}
```

#### List Groups
```http
GET /api/v1/groups?page=1&limit=20
Authorization: Bearer <access_token>
```

#### Add Member to Group
```http
POST /api/v1/groups/{id}/members
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "user_id": "user-uuid"
}
```

#### Remove Member from Group
```http
DELETE /api/v1/groups/{id}/members/{user_id}
Authorization: Bearer <access_token>
```

### Authorization Endpoints

#### Create Role
```http
POST /api/v1/roles
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "admin",
  "description": "Administrator role"
}
```

#### Create Permission
```http
POST /api/v1/permissions
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "resource": "users",
  "action": "read",
  "description": "Read user information"
}
```

#### Assign Permission to Role
```http
POST /api/v1/roles/{id}/permissions
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "permission_id": "permission-uuid"
}
```

#### Assign Role to User
```http
POST /api/v1/users/{id}/roles
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "role_id": "role-uuid"
}
```

#### Check Permission
```http
POST /api/v1/authz/check
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "user_id": "user-uuid",
  "resource": "users",
  "action": "read"
}
```

## 🛠️ SDK Usage

### Go SDK

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
    client := arasauth.NewClient("http://localhost:8080")
    
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
    
    // Check permission
    hasPermission, err := client.CheckPermission(context.Background(), user.ID, "users", "read")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Has permission: %v\n", hasPermission)
}
```

### Python SDK

```python
from aras_auth import AuthClient

# Initialize client
client = AuthClient("http://localhost:8080")

# Login
auth_response = client.login("user@example.com", "password")
print(f"Access token: {auth_response.access_token}")

# Get current user
user = client.get_current_user()
print(f"User: {user.first_name} {user.last_name}")

# Check permission
has_permission = client.check_permission(user.id, "users", "read")
print(f"Has permission: {has_permission}")

# Create group
group = client.create_group("Developers", "Development team")
print(f"Created group: {group.name}")

# Add member to group
client.add_member(group.id, user.id)
print("Member added to group")
```

## 🗄️ Database Schema

### Core Tables

- `users` - User accounts
- `groups` - User groups
- `user_groups` - Many-to-many relationship between users and groups
- `roles` - Roles
- `permissions` - Permissions
- `role_permissions` - Many-to-many relationship between roles and permissions
- `user_roles` - User role assignments
- `group_roles` - Group role assignments
- `refresh_tokens` - Refresh token storage
- `providers` - Identity provider registry

### Initial Data

The service comes with pre-seeded data:
- **Roles**: admin, user, moderator
- **Permissions**: users:create, users:read, users:update, users:delete, etc.
- **Role-Permission Assignments**: Admin has all permissions, user has basic read permissions

## 🔒 Security Considerations

### Production Deployment

1. **Change default secrets**: Update `JWT_SECRET_KEY` and admin credentials
2. **Use HTTPS**: Always use HTTPS in production
3. **Database security**: Use strong passwords and enable SSL
4. **Network security**: Use proper firewall rules
5. **Monitoring**: Implement logging and monitoring
6. **Backup**: Regular database backups

### JWT Security

- Access tokens expire in 15 minutes (configurable)
- Refresh tokens expire in 7 days (configurable)
- Refresh token rotation is implemented
- Tokens are signed with HS256

### Password Security

- Passwords are hashed with bcrypt (cost 12)
- Minimum password length is 8 characters
- Password validation can be extended

## 🧪 Testing

### Run Tests
```bash
go test -v ./...
```

### Run with Coverage
```bash
go test -v -cover ./...
```

### Integration Tests
```bash
# Start test database
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
go test -v -tags=integration ./...

# Cleanup
docker-compose -f docker-compose.test.yml down
```

## 🚀 Deployment

### Docker Deployment

1. **Build image**
   ```bash
   docker build -t aras-auth:latest .
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Scale service**
   ```bash
   docker-compose up -d --scale aras_auth=3
   ```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: aras-auth
spec:
  replicas: 3
  selector:
    matchLabels:
      app: aras-auth
  template:
    metadata:
      labels:
        app: aras-auth
    spec:
      containers:
      - name: aras-auth
        image: aras-auth:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres-service"
        - name: JWT_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: aras-auth-secrets
              key: jwt-secret
```

## 📈 Monitoring

### Health Check
```http
GET /health
```

### Metrics (Future Enhancement)
- Request count
- Response time
- Error rate
- Active users
- Token usage

### Logging
- Structured logging with zap
- Request/response logging
- Error logging
- Audit logging for security events

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

### Development Setup

1. **Clone repository**
   ```bash
   git clone https://github.com/aras-services/aras-auth.git
   cd aras-auth
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run in development mode**
   ```bash
   make dev
   ```

4. **Run tests**
   ```bash
   make test
   ```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: [GitHub Wiki](https://github.com/aras-services/aras-auth/wiki)
- **Issues**: [GitHub Issues](https://github.com/aras-services/aras-auth/issues)
- **Discussions**: [GitHub Discussions](https://github.com/aras-services/aras-auth/discussions)

## 🗺️ Roadmap

### Version 2.0
- [ ] LDAP provider implementation
- [ ] Multi-factor authentication (TOTP)
- [ ] OAuth2/OIDC server capabilities
- [ ] Advanced session management

### Version 3.0
- [ ] ABAC support (attribute-based policies)
- [ ] Policy decision engine
- [ ] Audit logging
- [ ] Advanced monitoring

### Version 4.0
- [ ] SAML provider
- [ ] SSO capabilities
- [ ] Advanced session management
- [ ] GraphQL API

## 🙏 Acknowledgments

- Inspired by the [ORY](https://www.ory.sh/) stack
- Built with [Go](https://golang.org/)
- Database: [PostgreSQL](https://www.postgresql.org/)
- JWT: [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- HTTP Router: [Chi](https://github.com/go-chi/chi)


