# ArasAuth Python SDK

A Python client library for the ArasAuth authentication and authorization service.

## Installation

```bash
pip install aras-auth
```

## Usage

### Basic Authentication

```python
from aras_auth import AuthClient

# Initialize the client
client = AuthClient("http://localhost:8080")

# Login
auth_response = client.login("user@example.com", "password")
print(f"Access token: {auth_response.access_token}")

# Get current user
user = client.get_current_user()
print(f"User: {user.email}")
```

### User Management

```python
# List users
users_response = client.list_users(page=1, limit=10)
print(f"Total users: {users_response.total}")

# Get specific user
user = client.get_user("user-id")

# Update user
updated_user = client.update_user("user-id", first_name="John", last_name="Doe")

# Delete user
client.delete_user("user-id")
```

### Group Management

```python
# Create group
group = client.create_group("Developers", "Development team")

# List groups
groups_response = client.list_groups()

# Add member to group
client.add_member("group-id", "user-id")

# Get group members
members = client.get_members("group-id")
```

### Role-Based Access Control

```python
# Create role
role = client.create_role("admin", "Administrator role")

# Create permission
permission = client.create_permission("users", "read", "Read user information")

# Assign permission to role
client.assign_permission_to_role("role-id", "permission-id")

# Assign role to user
client.assign_role_to_user("user-id", "role-id")

# Check permission
has_permission = client.check_permission("user-id", "users", "read")
```

### Error Handling

```python
try:
    user = client.get_user("invalid-id")
except Exception as e:
    print(f"Error: {e}")
```

## API Reference

### AuthClient

The main client class for interacting with the ArasAuth service.

#### Constructor

```python
AuthClient(base_url: str, timeout: int = 30)
```

- `base_url`: The base URL of the ArasAuth service
- `timeout`: Request timeout in seconds (default: 30)

#### Authentication Methods

- `login(email: str, password: str) -> AuthResponse`
- `register(email: str, password: str, first_name: str = '', last_name: str = '') -> User`
- `refresh_token(refresh_token: str) -> AuthResponse`
- `logout(refresh_token: str) -> None`
- `get_current_user() -> User`
- `check_permission(user_id: str, resource: str, action: str) -> bool`

#### User Management Methods

- `list_users(page: int = 1, limit: int = 20) -> ListResponse`
- `get_user(user_id: str) -> User`
- `update_user(user_id: str, **kwargs) -> User`
- `delete_user(user_id: str) -> None`

#### Group Management Methods

- `list_groups(page: int = 1, limit: int = 20) -> ListResponse`
- `create_group(name: str, description: str = '') -> Group`
- `get_group(group_id: str) -> Group`
- `update_group(group_id: str, **kwargs) -> Group`
- `delete_group(group_id: str) -> None`
- `add_member(group_id: str, user_id: str) -> None`
- `remove_member(group_id: str, user_id: str) -> None`
- `get_members(group_id: str) -> List[User]`

#### Role Management Methods

- `list_roles(page: int = 1, limit: int = 20) -> ListResponse`
- `create_role(name: str, description: str = '') -> Role`
- `get_role(role_id: str) -> Role`
- `update_role(role_id: str, **kwargs) -> Role`
- `delete_role(role_id: str) -> None`
- `assign_role_to_user(user_id: str, role_id: str) -> None`
- `remove_role_from_user(user_id: str, role_id: str) -> None`
- `get_user_roles(user_id: str) -> List[Role]`
- `assign_role_to_group(group_id: str, role_id: str) -> None`
- `remove_role_from_group(group_id: str, role_id: str) -> None`
- `get_group_roles(group_id: str) -> List[Role]`

#### Permission Management Methods

- `list_permissions(page: int = 1, limit: int = 20) -> ListResponse`
- `create_permission(resource: str, action: str, description: str = '') -> Permission`
- `get_permission(permission_id: str) -> Permission`
- `update_permission(permission_id: str, **kwargs) -> Permission`
- `delete_permission(permission_id: str) -> None`
- `assign_permission_to_role(role_id: str, permission_id: str) -> None`
- `remove_permission_from_role(role_id: str, permission_id: str) -> None`
- `get_role_permissions(role_id: str) -> List[Permission]`

## Models

### User

```python
@dataclass
class User:
    id: str
    email: str
    first_name: str
    last_name: str
    status: str
    email_verified: bool
    created_at: str
    updated_at: str
```

### Group

```python
@dataclass
class Group:
    id: str
    name: str
    description: str
    created_at: str
    updated_at: str
```

### Role

```python
@dataclass
class Role:
    id: str
    name: str
    description: str
    created_at: str
    updated_at: str
```

### Permission

```python
@dataclass
class Permission:
    id: str
    resource: str
    action: str
    description: str
    created_at: str
    updated_at: str
```

### AuthResponse

```python
@dataclass
class AuthResponse:
    access_token: str
    refresh_token: str
    expires_in: int
    token_type: str
    user: User
```

## License

MIT License


