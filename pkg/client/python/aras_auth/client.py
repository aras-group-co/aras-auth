"""
ArasAuth Python SDK Client

A Python client library for the ArasAuth authentication and authorization service.
"""

import json
import requests
from typing import Optional, Dict, Any, List
from urllib.parse import urljoin

from .models import User, Group, Role, Permission, AuthResponse, ListResponse


class AuthClient:
    """Client for interacting with the ArasAuth service"""
    
    def __init__(self, base_url: str, timeout: int = 30):
        """
        Initialize the AuthClient
        
        Args:
            base_url: The base URL of the ArasAuth service
            timeout: Request timeout in seconds
        """
        self.base_url = base_url.rstrip('/')
        self.timeout = timeout
        self.session = requests.Session()
        self.token = None
        
        # Set default headers
        self.session.headers.update({
            'Content-Type': 'application/json',
            'User-Agent': 'ArasAuth-Python-SDK/1.0.0'
        })
    
    def set_token(self, token: str) -> None:
        """Set the authentication token for the client"""
        self.token = token
        self.session.headers.update({'Authorization': f'Bearer {token}'})
    
    def _make_request(self, method: str, endpoint: str, data: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """
        Make an HTTP request to the API
        
        Args:
            method: HTTP method
            endpoint: API endpoint
            data: Request data
            
        Returns:
            Response data as dictionary
            
        Raises:
            requests.RequestException: If the request fails
        """
        url = urljoin(self.base_url, endpoint)
        
        try:
            response = self.session.request(
                method=method,
                url=url,
                json=data,
                timeout=self.timeout
            )
            response.raise_for_status()
            
            return response.json()
        except requests.exceptions.RequestException as e:
            raise Exception(f"Request failed: {str(e)}")
    
    def _handle_response(self, response_data: Dict[str, Any]) -> Any:
        """
        Handle API response and extract data
        
        Args:
            response_data: Response data from API
            
        Returns:
            Extracted data from response
        """
        if not response_data.get('success', False):
            error = response_data.get('error', 'Unknown error')
            message = response_data.get('message', '')
            raise Exception(f"API error: {error} - {message}")
        
        return response_data.get('data')
    
    # Authentication methods
    
    def login(self, email: str, password: str) -> AuthResponse:
        """
        Authenticate a user and return tokens
        
        Args:
            email: User email
            password: User password
            
        Returns:
            AuthResponse with tokens and user info
        """
        data = {
            'email': email,
            'password': password
        }
        
        response_data = self._make_request('POST', '/api/v1/auth/login', data)
        auth_data = self._handle_response(response_data)
        
        auth_response = AuthResponse.from_dict(auth_data)
        self.set_token(auth_response.access_token)
        
        return auth_response
    
    def register(self, email: str, password: str, first_name: str = '', last_name: str = '') -> User:
        """
        Register a new user
        
        Args:
            email: User email
            password: User password
            first_name: User first name
            last_name: User last name
            
        Returns:
            Created user
        """
        data = {
            'email': email,
            'password': password,
            'first_name': first_name,
            'last_name': last_name
        }
        
        response_data = self._make_request('POST', '/api/v1/auth/register', data)
        user_data = self._handle_response(response_data)
        
        return User.from_dict(user_data)
    
    def refresh_token(self, refresh_token: str) -> AuthResponse:
        """
        Refresh access token using refresh token
        
        Args:
            refresh_token: Refresh token
            
        Returns:
            AuthResponse with new tokens
        """
        data = {
            'refresh_token': refresh_token
        }
        
        response_data = self._make_request('POST', '/api/v1/auth/refresh', data)
        auth_data = self._handle_response(response_data)
        
        auth_response = AuthResponse.from_dict(auth_data)
        self.set_token(auth_response.access_token)
        
        return auth_response
    
    def logout(self, refresh_token: str) -> None:
        """
        Logout user by invalidating refresh token
        
        Args:
            refresh_token: Refresh token to invalidate
        """
        data = {
            'refresh_token': refresh_token
        }
        
        self._make_request('POST', '/api/v1/auth/logout', data)
        self.token = None
        self.session.headers.pop('Authorization', None)
    
    def get_current_user(self) -> User:
        """
        Get current authenticated user
        
        Returns:
            Current user
        """
        response_data = self._make_request('GET', '/api/v1/users/me')
        user_data = self._handle_response(response_data)
        
        return User.from_dict(user_data)
    
    def check_permission(self, user_id: str, resource: str, action: str) -> bool:
        """
        Check if a user has a specific permission
        
        Args:
            user_id: User ID
            resource: Resource name
            action: Action name
            
        Returns:
            True if user has permission, False otherwise
        """
        data = {
            'user_id': user_id,
            'resource': resource,
            'action': action
        }
        
        response_data = self._make_request('POST', '/api/v1/authz/check', data)
        permission_data = self._handle_response(response_data)
        
        return permission_data.get('has_permission', False)
    
    # User management methods
    
    def list_users(self, page: int = 1, limit: int = 20) -> ListResponse:
        """
        List users with pagination
        
        Args:
            page: Page number
            limit: Items per page
            
        Returns:
            ListResponse with users
        """
        params = f'?page={page}&limit={limit}'
        response_data = self._make_request('GET', f'/api/v1/users{params}')
        users_data = self._handle_response(response_data)
        
        return ListResponse.from_dict(users_data, User)
    
    def get_user(self, user_id: str) -> User:
        """
        Get a specific user by ID
        
        Args:
            user_id: User ID
            
        Returns:
            User
        """
        response_data = self._make_request('GET', f'/api/v1/users/{user_id}')
        user_data = self._handle_response(response_data)
        
        return User.from_dict(user_data)
    
    def update_user(self, user_id: str, **kwargs) -> User:
        """
        Update a user
        
        Args:
            user_id: User ID
            **kwargs: User fields to update
            
        Returns:
            Updated user
        """
        response_data = self._make_request('PUT', f'/api/v1/users/{user_id}', kwargs)
        user_data = self._handle_response(response_data)
        
        return User.from_dict(user_data)
    
    def delete_user(self, user_id: str) -> None:
        """
        Delete a user
        
        Args:
            user_id: User ID
        """
        self._make_request('DELETE', f'/api/v1/users/{user_id}')
    
    # Group management methods
    
    def list_groups(self, page: int = 1, limit: int = 20) -> ListResponse:
        """
        List groups with pagination
        
        Args:
            page: Page number
            limit: Items per page
            
        Returns:
            ListResponse with groups
        """
        params = f'?page={page}&limit={limit}'
        response_data = self._make_request('GET', f'/api/v1/groups{params}')
        groups_data = self._handle_response(response_data)
        
        return ListResponse.from_dict(groups_data, Group)
    
    def create_group(self, name: str, description: str = '') -> Group:
        """
        Create a new group
        
        Args:
            name: Group name
            description: Group description
            
        Returns:
            Created group
        """
        data = {
            'name': name,
            'description': description
        }
        
        response_data = self._make_request('POST', '/api/v1/groups', data)
        group_data = self._handle_response(response_data)
        
        return Group.from_dict(group_data)
    
    def get_group(self, group_id: str) -> Group:
        """
        Get a specific group by ID
        
        Args:
            group_id: Group ID
            
        Returns:
            Group
        """
        response_data = self._make_request('GET', f'/api/v1/groups/{group_id}')
        group_data = self._handle_response(response_data)
        
        return Group.from_dict(group_data)
    
    def update_group(self, group_id: str, **kwargs) -> Group:
        """
        Update a group
        
        Args:
            group_id: Group ID
            **kwargs: Group fields to update
            
        Returns:
            Updated group
        """
        response_data = self._make_request('PUT', f'/api/v1/groups/{group_id}', kwargs)
        group_data = self._handle_response(response_data)
        
        return Group.from_dict(group_data)
    
    def delete_group(self, group_id: str) -> None:
        """
        Delete a group
        
        Args:
            group_id: Group ID
        """
        self._make_request('DELETE', f'/api/v1/groups/{group_id}')
    
    def add_member(self, group_id: str, user_id: str) -> None:
        """
        Add a user to a group
        
        Args:
            group_id: Group ID
            user_id: User ID
        """
        data = {'user_id': user_id}
        self._make_request('POST', f'/api/v1/groups/{group_id}/members', data)
    
    def remove_member(self, group_id: str, user_id: str) -> None:
        """
        Remove a user from a group
        
        Args:
            group_id: Group ID
            user_id: User ID
        """
        self._make_request('DELETE', f'/api/v1/groups/{group_id}/members/{user_id}')
    
    def get_members(self, group_id: str) -> List[User]:
        """
        Get all members of a group
        
        Args:
            group_id: Group ID
            
        Returns:
            List of users
        """
        response_data = self._make_request('GET', f'/api/v1/groups/{group_id}/members')
        members_data = self._handle_response(response_data)
        
        return [User.from_dict(member) for member in members_data]
    
    # Role management methods
    
    def list_roles(self, page: int = 1, limit: int = 20) -> ListResponse:
        """
        List roles with pagination
        
        Args:
            page: Page number
            limit: Items per page
            
        Returns:
            ListResponse with roles
        """
        params = f'?page={page}&limit={limit}'
        response_data = self._make_request('GET', f'/api/v1/roles{params}')
        roles_data = self._handle_response(response_data)
        
        return ListResponse.from_dict(roles_data, Role)
    
    def create_role(self, name: str, description: str = '') -> Role:
        """
        Create a new role
        
        Args:
            name: Role name
            description: Role description
            
        Returns:
            Created role
        """
        data = {
            'name': name,
            'description': description
        }
        
        response_data = self._make_request('POST', '/api/v1/roles', data)
        role_data = self._handle_response(response_data)
        
        return Role.from_dict(role_data)
    
    def get_role(self, role_id: str) -> Role:
        """
        Get a specific role by ID
        
        Args:
            role_id: Role ID
            
        Returns:
            Role
        """
        response_data = self._make_request('GET', f'/api/v1/roles/{role_id}')
        role_data = self._handle_response(response_data)
        
        return Role.from_dict(role_data)
    
    def update_role(self, role_id: str, **kwargs) -> Role:
        """
        Update a role
        
        Args:
            role_id: Role ID
            **kwargs: Role fields to update
            
        Returns:
            Updated role
        """
        response_data = self._make_request('PUT', f'/api/v1/roles/{role_id}', kwargs)
        role_data = self._handle_response(response_data)
        
        return Role.from_dict(role_data)
    
    def delete_role(self, role_id: str) -> None:
        """
        Delete a role
        
        Args:
            role_id: Role ID
        """
        self._make_request('DELETE', f'/api/v1/roles/{role_id}')
    
    def assign_role_to_user(self, user_id: str, role_id: str) -> None:
        """
        Assign a role to a user
        
        Args:
            user_id: User ID
            role_id: Role ID
        """
        data = {'role_id': role_id}
        self._make_request('POST', f'/api/v1/users/{user_id}/roles', data)
    
    def remove_role_from_user(self, user_id: str, role_id: str) -> None:
        """
        Remove a role from a user
        
        Args:
            user_id: User ID
            role_id: Role ID
        """
        self._make_request('DELETE', f'/api/v1/users/{user_id}/roles/{role_id}')
    
    def get_user_roles(self, user_id: str) -> List[Role]:
        """
        Get all roles for a user
        
        Args:
            user_id: User ID
            
        Returns:
            List of roles
        """
        response_data = self._make_request('GET', f'/api/v1/users/{user_id}/roles')
        roles_data = self._handle_response(response_data)
        
        return [Role.from_dict(role) for role in roles_data]
    
    def assign_role_to_group(self, group_id: str, role_id: str) -> None:
        """
        Assign a role to a group
        
        Args:
            group_id: Group ID
            role_id: Role ID
        """
        data = {'role_id': role_id}
        self._make_request('POST', f'/api/v1/groups/{group_id}/roles', data)
    
    def remove_role_from_group(self, group_id: str, role_id: str) -> None:
        """
        Remove a role from a group
        
        Args:
            group_id: Group ID
            role_id: Role ID
        """
        self._make_request('DELETE', f'/api/v1/groups/{group_id}/roles/{role_id}')
    
    def get_group_roles(self, group_id: str) -> List[Role]:
        """
        Get all roles for a group
        
        Args:
            group_id: Group ID
            
        Returns:
            List of roles
        """
        response_data = self._make_request('GET', f'/api/v1/groups/{group_id}/roles')
        roles_data = self._handle_response(response_data)
        
        return [Role.from_dict(role) for role in roles_data]
    
    # Permission management methods
    
    def list_permissions(self, page: int = 1, limit: int = 20) -> ListResponse:
        """
        List permissions with pagination
        
        Args:
            page: Page number
            limit: Items per page
            
        Returns:
            ListResponse with permissions
        """
        params = f'?page={page}&limit={limit}'
        response_data = self._make_request('GET', f'/api/v1/permissions{params}')
        permissions_data = self._handle_response(response_data)
        
        return ListResponse.from_dict(permissions_data, Permission)
    
    def create_permission(self, resource: str, action: str, description: str = '') -> Permission:
        """
        Create a new permission
        
        Args:
            resource: Resource name
            action: Action name
            description: Permission description
            
        Returns:
            Created permission
        """
        data = {
            'resource': resource,
            'action': action,
            'description': description
        }
        
        response_data = self._make_request('POST', '/api/v1/permissions', data)
        permission_data = self._handle_response(response_data)
        
        return Permission.from_dict(permission_data)
    
    def get_permission(self, permission_id: str) -> Permission:
        """
        Get a specific permission by ID
        
        Args:
            permission_id: Permission ID
            
        Returns:
            Permission
        """
        response_data = self._make_request('GET', f'/api/v1/permissions/{permission_id}')
        permission_data = self._handle_response(response_data)
        
        return Permission.from_dict(permission_data)
    
    def update_permission(self, permission_id: str, **kwargs) -> Permission:
        """
        Update a permission
        
        Args:
            permission_id: Permission ID
            **kwargs: Permission fields to update
            
        Returns:
            Updated permission
        """
        response_data = self._make_request('PUT', f'/api/v1/permissions/{permission_id}', kwargs)
        permission_data = self._handle_response(response_data)
        
        return Permission.from_dict(permission_data)
    
    def delete_permission(self, permission_id: str) -> None:
        """
        Delete a permission
        
        Args:
            permission_id: Permission ID
        """
        self._make_request('DELETE', f'/api/v1/permissions/{permission_id}')
    
    def assign_permission_to_role(self, role_id: str, permission_id: str) -> None:
        """
        Assign a permission to a role
        
        Args:
            role_id: Role ID
            permission_id: Permission ID
        """
        data = {'permission_id': permission_id}
        self._make_request('POST', f'/api/v1/roles/{role_id}/permissions', data)
    
    def remove_permission_from_role(self, role_id: str, permission_id: str) -> None:
        """
        Remove a permission from a role
        
        Args:
            role_id: Role ID
            permission_id: Permission ID
        """
        self._make_request('DELETE', f'/api/v1/roles/{role_id}/permissions/{permission_id}')
    
    def get_role_permissions(self, role_id: str) -> List[Permission]:
        """
        Get all permissions for a role
        
        Args:
            role_id: Role ID
            
        Returns:
            List of permissions
        """
        response_data = self._make_request('GET', f'/api/v1/roles/{role_id}/permissions')
        permissions_data = self._handle_response(response_data)
        
        return [Permission.from_dict(permission) for permission in permissions_data]


